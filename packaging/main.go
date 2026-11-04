// DeviceBridge's build-host package tool. It is not installed on devices.
package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"debug/elf"
	"debug/macho"
	_ "embed"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

//go:embed install.sh
var installer string

type options struct {
	component, version, target, binary, revision, output string
}

var binaries = map[string]string{
	"agent": "devicebridge-agent", "control": "devicebridge-control",
	"relay": "devicebridge-relay", "cli": "devicebridge",
}

var versionPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9.-]{0,63}$`)
var revisionPattern = regexp.MustCompile(`^[a-f0-9]{40}$`)

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

func run(args []string, stdout, stderr io.Writer) int {
	var o options
	f := flag.NewFlagSet("devicebridge-package", flag.ContinueOnError)
	f.SetOutput(stderr)
	f.StringVar(&o.component, "component", "", "agent, control, relay, or cli")
	f.StringVar(&o.version, "version", "", "unique package version, e.g. 0.1.0-preview.1")
	f.StringVar(&o.target, "target", "", "linux-armv7, linux-arm64, linux-amd64, darwin-arm64, darwin-amd64")
	f.StringVar(&o.binary, "binary", "", "explicit prebuilt executable (never built or executed by this tool)")
	f.StringVar(&o.revision, "revision", "", "declared source commit, 40 lowercase hex characters")
	f.StringVar(&o.output, "output", "dist", "output directory; existing archives are never overwritten")
	if err := f.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if f.NArg() != 0 {
		fmt.Fprintln(stderr, "unexpected positional arguments")
		return 2
	}
	path, err := build(o)
	if err != nil {
		fmt.Fprintln(stderr, "package:", err)
		return 1
	}
	if _, err := fmt.Fprintln(stdout, path); err != nil {
		return 1
	}
	return 0
}

func (o options) validate() error {
	if _, ok := binaries[o.component]; !ok {
		return errors.New("unsupported component")
	}
	if !versionPattern.MatchString(o.version) || strings.Contains(o.version, "..") {
		return errors.New("version must be 1..64 ASCII letters/digits/dots/hyphens, without '..', starting with a letter/digit")
	}
	switch o.target {
	case "linux-armv7", "linux-arm64", "linux-amd64", "darwin-arm64", "darwin-amd64":
	default:
		return errors.New("unsupported target")
	}
	if !revisionPattern.MatchString(o.revision) {
		return errors.New("revision must be a full lowercase 40-character source commit")
	}
	if o.binary == "" || o.output == "" {
		return errors.New("binary and output paths are required")
	}
	return nil
}

func checkExecutable(f *os.File, target string) error {
	info, err := f.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || info.Mode()&0111 == 0 {
		return errors.New("input must be a regular executable file")
	}
	if strings.HasPrefix(target, "linux-") {
		e, err := elf.NewFile(f)
		if err != nil {
			return errors.New("Linux package requires an ELF executable")
		}
		want := map[string]elf.Machine{"linux-armv7": elf.EM_ARM, "linux-arm64": elf.EM_AARCH64, "linux-amd64": elf.EM_X86_64}[target]
		class := elf.ELFCLASS64
		if target == "linux-armv7" {
			class = elf.ELFCLASS32
		}
		if e.Machine != want || e.Class != class || e.Data != elf.ELFDATA2LSB || (e.Type != elf.ET_EXEC && e.Type != elf.ET_DYN) {
			return errors.New("ELF architecture/class/type does not match target")
		}
	} else {
		m, err := macho.NewFile(f)
		if err != nil {
			return errors.New("macOS package requires a single-architecture Mach-O executable")
		}
		want := macho.CpuArm64
		if target == "darwin-amd64" {
			want = macho.CpuAmd64
		}
		if m.Cpu != want || m.Type != macho.TypeExec {
			return errors.New("Mach-O architecture/type does not match target")
		}
	}
	return nil
}

func build(o options) (result string, err error) {
	if err = o.validate(); err != nil {
		return "", err
	}
	info, err := os.Lstat(o.binary)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", errors.New("binary must be a regular file, not a symlink")
	}
	f, err := os.Open(o.binary)
	if err != nil {
		return "", err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil {
		return "", err
	}
	if !os.SameFile(info, opened) {
		return "", errors.New("input changed while opening")
	}
	if err = checkExecutable(f, o.target); err != nil {
		return "", err
	}
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}
	digest := hex.EncodeToString(h.Sum(nil))
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	if err = os.MkdirAll(o.output, 0755); err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp(o.output, ".devicebridge-package-*")
	if err != nil {
		return "", err
	}
	defer func() { tmp.Close(); os.Remove(tmp.Name()) }()
	name := "devicebridge-" + o.component + "-" + o.version + "-" + o.target
	result = filepath.Join(o.output, name+".tar.gz")
	gz := gzip.NewWriter(tmp)
	tw := tar.NewWriter(gz)
	add := func(path string, mode int64, data string) error {
		if err := tw.WriteHeader(&tar.Header{Name: name + "/" + path, Mode: mode, Size: int64(len(data)), ModTime: time.Unix(0, 0)}); err != nil {
			return err
		}
		_, err := io.WriteString(tw, data)
		return err
	}
	script := strings.NewReplacer("@COMPONENT@", o.component, "@VERSION@", o.version, "@TARGET@", o.target, "@BINARY@", binaries[o.component], "@SHA256@", digest).Replace(installer)
	readme := fmt.Sprintf("DeviceBridge M01 preview: %s %s (%s)\nDeclared source: %s\nPayload SHA-256: %s\n\nInstall: sh install.sh\nHelp: sh install.sh --help\nDefault prefix: $HOME/.local/devicebridge/%s\nStatus/rollback/uninstall: sh PREFIX/manage.sh --status|--rollback|--uninstall\nRun: PREFIX/bin/%s --help\n\nFoundation only: no remote SSH/tunnel/enrollment. Do not run as a production service.\nUse trusted archives; SHA-256 detects corruption, not publisher identity.\nLinux ARM ISA, libc/loader and distribution compatibility require separate validation.\nInstallation never starts a service, changes PATH or writes configuration/secrets.\nNo compiler, Docker, Python or Node is needed on the destination.\n", o.component, o.version, o.target, o.revision, digest, o.component, binaries[o.component])
	for _, entry := range []struct {
		name, data string
		mode       int64
	}{
		{"install.sh", script, 0755}, {"README.txt", readme, 0644},
		{"SHA256SUMS", digest + "  bin/" + binaries[o.component] + "\n", 0644},
	} {
		if err = add(entry.name, entry.mode, entry.data); err != nil {
			return "", err
		}
	}
	if err = tw.WriteHeader(&tar.Header{Name: name + "/bin/" + binaries[o.component], Mode: 0755, Size: opened.Size(), ModTime: time.Unix(0, 0)}); err != nil {
		return "", err
	}
	copyHash := sha256.New()
	if _, err = io.Copy(tw, io.TeeReader(f, copyHash)); err != nil {
		return "", err
	}
	if hex.EncodeToString(copyHash.Sum(nil)) != digest {
		return "", errors.New("input changed while packaging")
	}
	if err = tw.Close(); err != nil {
		return "", err
	}
	if err = gz.Close(); err != nil {
		return "", err
	}
	if err = tmp.Close(); err != nil {
		return "", err
	}
	// Link publishes a complete archive without replacing an existing file/symlink.
	if err = os.Chmod(tmp.Name(), 0644); err != nil {
		return "", err
	}
	if err = os.Link(tmp.Name(), result); err != nil {
		return "", err
	}
	published := result
	defer func() {
		if err != nil {
			os.Remove(published)
		}
	}()
	archive, err := os.Open(result)
	if err != nil {
		return "", err
	}
	h.Reset()
	_, err = io.Copy(h, archive)
	closeErr := archive.Close()
	if err != nil {
		return "", err
	}
	if closeErr != nil {
		return "", closeErr
	}
	sum, err := os.OpenFile(result+".sha256", os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return "", err
	}
	_, err = fmt.Fprintf(sum, "%x  %s\n", h.Sum(nil), filepath.Base(result))
	closeErr = sum.Close()
	if err != nil {
		os.Remove(result + ".sha256")
		return "", err
	}
	if closeErr != nil {
		os.Remove(result + ".sha256")
		return "", closeErr
	}
	return result, nil
}
