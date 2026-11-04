package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const testRevision = "0123456789abcdef0123456789abcdef01234567"

func temp(t *testing.T) string {
	t.Helper()
	p, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func testOptions(t *testing.T) options {
	t.Helper()
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	return options{component: "cli", version: "0.1.0-test.1", target: runtime.GOOS + "-" + runtime.GOARCH, binary: binary, revision: testRevision, output: temp(t)}
}

func unpack(t *testing.T, o options) string {
	t.Helper()
	p, err := build(o)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer gz.Close()
	r := tar.NewReader(gz)
	dest := temp(t)
	root := ""
	for {
		h, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if h.Typeflag != tar.TypeReg || filepath.IsAbs(h.Name) || strings.Contains(h.Name, "..") {
			t.Fatalf("unsafe generated entry: %+v", h)
		}
		if root == "" {
			root = strings.Split(h.Name, "/")[0]
		}
		out := filepath.Join(dest, h.Name)
		if err := os.MkdirAll(filepath.Dir(out), 0700); err != nil {
			t.Fatal(err)
		}
		data, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(out, data, os.FileMode(h.Mode)); err != nil {
			t.Fatal(err)
		}
	}
	return filepath.Join(dest, root)
}

func shell(t *testing.T, script, home string, wantSuccess bool, args ...string) string {
	t.Helper()
	cmd := exec.Command("sh", append([]string{script}, args...)...)
	cmd.Env = append(os.Environ(), "HOME="+home)
	output, err := cmd.CombinedOutput()
	if (err == nil) != wantSuccess {
		t.Fatalf("sh %s %v: %v\n%s", script, args, err, output)
	}
	return string(output)
}

func readLink(t *testing.T, p string) string {
	t.Helper()
	value, err := os.Readlink(p)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func absent(t *testing.T, p string) {
	t.Helper()
	if _, err := os.Lstat(p); !os.IsNotExist(err) {
		t.Fatalf("expected absent %s: %v", p, err)
	}
}

func TestInstallLifecycle(t *testing.T) {
	o := testOptions(t)
	root := unpack(t, o)
	home := temp(t)
	prefix := filepath.Join(home, ".local/devicebridge/cli")
	shell(t, filepath.Join(root, "install.sh"), home, true)
	active := filepath.Join(prefix, "bin/devicebridge")
	first := readLink(t, active)
	if !strings.Contains(first, o.version) {
		t.Fatalf("active=%s", first)
	}
	data, _ := os.ReadFile(active)
	source, _ := os.ReadFile(o.binary)
	if !bytes.Equal(data, source) {
		t.Fatal("installed binary differs")
	}
	before, _ := os.Stat(active)
	shell(t, filepath.Join(root, "install.sh"), home, true)
	after, _ := os.Stat(active)
	if !os.SameFile(before, after) {
		t.Fatal("idempotent install rewrote payload")
	}
	absent(t, filepath.Join(prefix, "previous"))
	o.version = "0.1.0-test.2"
	upgraded := unpack(t, o)
	shell(t, filepath.Join(upgraded, "install.sh"), home, true)
	second := readLink(t, active)
	if first == second {
		t.Fatal("upgrade did not activate")
	}
	manager := filepath.Join(prefix, "manage.sh")
	status := shell(t, manager, home, true, "--status")
	if !strings.Contains(status, o.version) {
		t.Fatal(status)
	}
	shell(t, manager, home, true, "--rollback")
	if readLink(t, active) != first {
		t.Fatal("rollback failed")
	}
	shell(t, manager, home, true, "--rollback")
	if readLink(t, active) != second {
		t.Fatal("second rollback failed")
	}
	config := filepath.Join(prefix, "my-config.json")
	if err := os.WriteFile(config, []byte("user owned"), 0600); err != nil {
		t.Fatal(err)
	}
	shell(t, manager, home, true, "--uninstall")
	absent(t, active)
	absent(t, filepath.Join(prefix, "releases"))
	absent(t, manager)
	if data, _ := os.ReadFile(config); string(data) != "user owned" {
		t.Fatal("user config changed")
	}
}

func TestExplicitPrefixSpacesAndCleanUninstall(t *testing.T) {
	o := testOptions(t)
	root := unpack(t, o)
	home := temp(t)
	prefix := filepath.Join(home, "custom prefix")
	shell(t, filepath.Join(root, "install.sh"), home, true, "--prefix", prefix)
	shell(t, filepath.Join(prefix, "manage.sh"), home, false) // Cannot reinstall from installed manager.
	shell(t, filepath.Join(prefix, "manage.sh"), home, false, "--rollback")
	shell(t, filepath.Join(prefix, "manage.sh"), home, true, "--uninstall")
	absent(t, prefix)
}

func TestRejectsInvalidPackagesBeforeWriting(t *testing.T) {
	for _, name := range []string{"corruption", "missing", "symlink", "symlink-bin", "target", "invalid-version"} {
		t.Run(name, func(t *testing.T) {
			o := testOptions(t)
			root := unpack(t, o)
			home := temp(t)
			payload := filepath.Join(root, "bin/devicebridge")
			script := filepath.Join(root, "install.sh")
			switch name {
			case "corruption":
				f, err := os.OpenFile(payload, os.O_APPEND|os.O_WRONLY, 0)
				if err != nil {
					t.Fatal(err)
				}
				_, _ = f.WriteString("corruption")
				_ = f.Close()
			case "missing":
				_ = os.Remove(payload)
			case "symlink":
				_ = os.Remove(payload)
				if err := os.Symlink(o.binary, payload); err != nil {
					t.Fatal(err)
				}
			case "symlink-bin":
				external := temp(t)
				data, _ := os.ReadFile(o.binary)
				_ = os.WriteFile(filepath.Join(external, "devicebridge"), data, 0700)
				_ = os.RemoveAll(filepath.Join(root, "bin"))
				_ = os.Symlink(external, filepath.Join(root, "bin"))
			case "target", "invalid-version":
				data, _ := os.ReadFile(script)
				if name == "target" {
					wrong := "linux-amd64"
					if o.target == wrong {
						wrong = "linux-arm64"
					}
					data = bytes.Replace(data, []byte("target='"+o.target+"'"), []byte("target='"+wrong+"'"), 1)
				} else {
					data = bytes.Replace(data, []byte("version='"+o.version+"'"), []byte("version='../escape'"), 1)
				}
				_ = os.WriteFile(script, data, 0700)
			}
			shell(t, script, home, false)
			absent(t, filepath.Join(home, ".local"))
		})
	}
}

func TestRejectsUnmanagedAndSymlinkDestinations(t *testing.T) {
	for _, name := range []string{"unmanaged", "symlink", "parent-symlink", "relative", "dotdot", "empty"} {
		t.Run(name, func(t *testing.T) {
			o := testOptions(t)
			root := unpack(t, o)
			home := temp(t)
			prefix := filepath.Join(home, "destination")
			switch name {
			case "unmanaged":
				_ = os.Mkdir(prefix, 0700)
				_ = os.WriteFile(filepath.Join(prefix, "keep"), []byte("keep"), 0600)
			case "symlink":
				_ = os.Symlink(temp(t), prefix)
			case "parent-symlink":
				_ = os.Symlink(temp(t), prefix)
				prefix = filepath.Join(prefix, "child")
			case "relative":
				prefix = "relative"
			case "dotdot":
				prefix = home + "/../escape"
			case "empty":
				prefix = ""
			}
			shell(t, filepath.Join(root, "install.sh"), home, false, "--prefix", prefix)
			if name == "unmanaged" {
				if data, _ := os.ReadFile(filepath.Join(prefix, "keep")); string(data) != "keep" {
					t.Fatal("unrelated file changed")
				}
			}
			absent(t, filepath.Join(home, ".local"))
		})
	}
}

func TestManagedStateFailuresLeaveInstallation(t *testing.T) {
	for _, name := range []string{"lock", "unrelated-bin", "foreign-link", "extra-release-file", "corrupt-previous", "version-reuse"} {
		t.Run(name, func(t *testing.T) {
			o := testOptions(t)
			root := unpack(t, o)
			home := temp(t)
			prefix := filepath.Join(home, "installed")
			script := filepath.Join(root, "install.sh")
			shell(t, script, home, true, "--prefix", prefix)
			active := filepath.Join(prefix, "bin/devicebridge")
			original := readLink(t, active)
			action := []string{"--prefix", prefix}
			switch name {
			case "lock":
				_ = os.Mkdir(filepath.Join(prefix, ".install-lock"), 0700)
			case "unrelated-bin":
				_ = os.Remove(active)
				_ = os.WriteFile(active, []byte("keep"), 0700)
			case "foreign-link":
				_ = os.Remove(active)
				_ = os.Symlink("../../elsewhere", active)
			case "extra-release-file":
				_ = os.WriteFile(filepath.Join(prefix, original[3:], "../extra"), []byte("keep"), 0600)
				action = append(action, "--uninstall")
			case "corrupt-previous":
				_ = os.Symlink("releases/../../elsewhere/devicebridge", filepath.Join(prefix, "previous"))
			case "version-reuse":
				copyPath := filepath.Join(temp(t), "modified")
				data, _ := os.ReadFile(o.binary)
				data = append(data, []byte("new payload")...)
				_ = os.WriteFile(copyPath, data, 0700)
				o.binary = copyPath
				o.output = temp(t)
				script = filepath.Join(unpack(t, o), "install.sh")
			}
			shell(t, script, home, false, action...)
			if _, err := os.Stat(filepath.Join(prefix, ".devicebridge-owner")); err != nil {
				t.Fatal("marker removed on failure")
			}
			if name == "unrelated-bin" {
				if data, _ := os.ReadFile(active); string(data) != "keep" {
					t.Fatal("unrelated binary changed")
				}
			}
			if name == "lock" {
				if _, err := os.Stat(filepath.Join(prefix, ".install-lock")); err != nil {
					t.Fatal("foreign lock removed")
				}
			}
			if name == "version-reuse" && readLink(t, active) != original {
				t.Fatal("version conflict changed active link")
			}
		})
	}
}

func TestPackageValidationAndNoOverwrite(t *testing.T) {
	base := testOptions(t)
	for _, mutate := range []func(*options){
		func(o *options) { o.component = "unknown" }, func(o *options) { o.version = "../bad" },
		func(o *options) { o.version = strings.Repeat("a", 65) }, func(o *options) { o.version = "" },
		func(o *options) { o.target = "unsupported" }, func(o *options) { o.revision = "short" },
		func(o *options) { o.binary = "" }, func(o *options) { o.output = "" },
	} {
		o := base
		mutate(&o)
		if _, err := build(o); err == nil {
			t.Fatalf("accepted invalid options: %+v", o)
		}
	}
	script := filepath.Join(temp(t), "script")
	_ = os.WriteFile(script, []byte("#!/bin/sh\nexit 0\n"), 0700)
	o := base
	o.binary = script
	if _, err := build(o); err == nil {
		t.Fatal("accepted non-native binary format")
	}
	link := filepath.Join(temp(t), "link")
	_ = os.Symlink(base.binary, link)
	o.binary = link
	if _, err := build(o); err == nil {
		t.Fatal("accepted binary symlink")
	}
	o = base
	if o.target == "darwin-arm64" {
		o.target = "darwin-amd64"
	} else if o.target == "linux-amd64" {
		o.target = "linux-arm64"
	} else {
		o.target = "linux-amd64"
	}
	if _, err := build(o); err == nil {
		t.Fatal("accepted mismatched architecture")
	}
	archive, err := build(base)
	if err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(archive)
	if _, err := build(base); err == nil {
		t.Fatal("overwrote existing package")
	}
	after, _ := os.ReadFile(archive)
	if !bytes.Equal(before, after) {
		t.Fatal("existing archive modified")
	}
	sum, _ := os.ReadFile(archive + ".sha256")
	if !strings.HasPrefix(string(sum), fmt.Sprintf("%x  ", sha256.Sum256(before))) {
		t.Fatal("bad archive checksum")
	}
	// A preexisting checksum path must not cause a partial archive or be removed.
	o = base
	o.version = "0.1.0-test.conflict"
	conflict := filepath.Join(o.output, "devicebridge-cli-"+o.version+"-"+o.target+".tar.gz")
	_ = os.Mkdir(conflict+".sha256", 0700)
	if _, err := build(o); err == nil {
		t.Fatal("accepted checksum path collision")
	}
	absent(t, conflict)
	if _, err := os.Stat(conflict + ".sha256"); err != nil {
		t.Fatal("unrelated checksum path removed")
	}
}

func TestArgumentsAndHelp(t *testing.T) {
	for _, args := range [][]string{{"--unknown"}, {"positional"}, {"--version"}, nil} {
		var out, err bytes.Buffer
		if run(args, &out, &err) == 0 {
			t.Fatalf("accepted %v", args)
		}
	}
	var out, err bytes.Buffer
	if run([]string{"--help"}, &out, &err) != 0 {
		t.Fatal("help failed")
	}
	root := unpack(t, testOptions(t))
	home := temp(t)
	for _, args := range [][]string{{"--wat"}, {"--prefix"}, {"--prefix", home, "--prefix", home}, {"--status", "--uninstall"}} {
		shell(t, filepath.Join(root, "install.sh"), home, false, args...)
	}
	shell(t, filepath.Join(root, "install.sh"), home, true, "--help")
	absent(t, filepath.Join(home, ".local"))
}

func TestInstallerHashUtilities(t *testing.T) {
	for _, selected := range []string{"sha256sum", "shasum", "openssl", "none"} {
		t.Run(selected, func(t *testing.T) {
			if selected != "none" {
				if _, err := exec.LookPath(selected); err != nil {
					t.Skip("hash utility unavailable")
				}
			}
			o := testOptions(t)
			root := unpack(t, o)
			home := temp(t)
			utilities := temp(t)
			for _, name := range []string{"uname", "dirname", "mkdir", "cp", "mv", "ln", "readlink", "chmod", "rm", "rmdir", "cat", "wc", "mktemp", selected} {
				if name == "none" {
					continue
				}
				path, err := exec.LookPath(name)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(path, filepath.Join(utilities, name)); err != nil {
					t.Fatal(err)
				}
			}
			cmd := exec.Command("sh", filepath.Join(root, "install.sh"))
			cmd.Env = append(os.Environ(), "HOME="+home, "PATH="+utilities)
			if selected == "none" {
				probe := exec.Command("sh", "-c", "command -v sha256sum || command -v shasum || command -v openssl")
				probe.Env = cmd.Env
				if probe.Run() == nil {
					t.Skip("shell exposes SHA utility as a built-in applet; PATH cannot simulate its absence")
				}
			}
			output, err := cmd.CombinedOutput()
			if selected == "none" {
				if err == nil {
					t.Fatal("installed without SHA-256 utility")
				}
				absent(t, filepath.Join(home, ".local"))
			} else if err != nil {
				t.Fatalf("hash utility %s: %v\n%s", selected, err, output)
			}
		})
	}
}

func TestFailedActivationPreservesRollback(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("permission failure requires an unprivileged UID")
	}
	for _, failure := range []string{"binary-directory", "history-rename"} {
		t.Run(failure, func(t *testing.T) {
			o := testOptions(t)
			home := temp(t)
			prefix := filepath.Join(home, "install")
			for _, version := range []string{"1", "2"} {
				o.version = version
				root := unpack(t, o)
				shell(t, filepath.Join(root, "install.sh"), home, true, "--prefix", prefix)
			}
			activePath := filepath.Join(prefix, "bin/devicebridge")
			beforeActive, beforePrevious := readLink(t, activePath), readLink(t, filepath.Join(prefix, "previous"))
			o.version = "3"
			root := unpack(t, o)
			cmd := exec.Command("sh", filepath.Join(root, "install.sh"), "--prefix", prefix)
			cmd.Env = append(os.Environ(), "HOME="+home)
			if failure == "binary-directory" {
				if err := os.Chmod(filepath.Join(prefix, "bin"), 0500); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { _ = os.Chmod(filepath.Join(prefix, "bin"), 0700) })
			} else {
				// Fail only the history commit; allow the executable commit and
				// restoration to use the real mv, reproducing the second failure path.
				mv, err := exec.LookPath("mv")
				if err != nil {
					t.Fatal(err)
				}
				shimDir := temp(t)
				shim := "#!/bin/sh\ncase \"$3\" in */previous) exit 1 ;; esac\nexec \"$DB011_REAL_MV\" \"$@\"\n"
				if err := os.WriteFile(filepath.Join(shimDir, "mv"), []byte(shim), 0700); err != nil {
					t.Fatal(err)
				}
				cmd.Env = append(cmd.Env, "PATH="+shimDir+":"+os.Getenv("PATH"), "DB011_REAL_MV="+mv)
			}
			output, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("activation unexpectedly succeeded: %s", output)
			}
			if failure == "binary-directory" {
				_ = os.Chmod(filepath.Join(prefix, "bin"), 0700)
			}
			if readLink(t, activePath) != beforeActive || readLink(t, filepath.Join(prefix, "previous")) != beforePrevious {
				t.Fatalf("failed update changed active/rollback state: %s", output)
			}
			shell(t, filepath.Join(prefix, "manage.sh"), home, true, "--rollback")
			if !strings.Contains(readLink(t, activePath), "/1-") {
				t.Fatal("rollback target was lost")
			}
		})
	}
}
