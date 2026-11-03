package control

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func testLogger() *slog.Logger { return slog.New(slog.NewJSONHandler(io.Discard, nil)) }

func TestProbeRoutes(t *testing.T) {
	c, _ := configFrom(nil)
	s := newService(c, testLogger())
	for _, test := range []struct {
		method, path, body string
		ready              bool
		status             int
	}{
		{"GET", "/healthz", "{\"status\":\"ok\"}\n", false, 200},
		{"HEAD", "/healthz", "", true, 200},
		{"GET", "/readyz", "{\"status\":\"ready\"}\n", true, 200},
		{"GET", "/readyz", "{\"status\":\"shutting_down\"}\n", false, 503},
		{"HEAD", "/readyz", "", false, 503},
		{"POST", "/healthz", "{\"error\":\"method_not_allowed\"}\n", true, 405},
		{"DELETE", "/readyz", "{\"error\":\"method_not_allowed\"}\n", true, 405},
		{"OPTIONS", "/healthz", "{\"error\":\"method_not_allowed\"}\n", true, 405},
		{"GET", "/", "{\"error\":\"not_found\"}\n", true, 404},
		{"GET", "/debug/pprof/", "{\"error\":\"not_found\"}\n", true, 404},
		{"GET", "/healthz/", "{\"error\":\"not_found\"}\n", true, 404},
		{"GET", "/x/../healthz", "{\"error\":\"not_found\"}\n", true, 404},
	} {
		t.Run(test.method+test.path+test.body, func(t *testing.T) {
			s.ready.Store(test.ready)
			w := httptest.NewRecorder()
			s.handle(w, httptest.NewRequest(test.method, test.path, nil))
			if w.Code != test.status || w.Body.String() != test.body {
				t.Fatalf("status=%d body=%q", w.Code, w.Body.String())
			}
			if w.Header().Get("Content-Type") != "application/json" || w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("X-Content-Type-Options") != "nosniff" {
				t.Fatalf("missing probe headers: %v", w.Header())
			}
			if test.status == 405 && w.Header().Get("Allow") != "GET, HEAD" {
				t.Fatal("missing Allow header")
			}
		})
	}
}

func startTestService(t *testing.T, c Config, s *service) (string, context.CancelFunc, <-chan error) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- s.serve(ctx, l, c, testLogger()) }()
	t.Cleanup(func() { cancel(); _ = s.http.Close() })
	return l.Addr().String(), cancel, done
}

func waitDone(t *testing.T, done <-chan error) error {
	t.Helper()
	select {
	case err := <-done:
		return err
	case <-time.After(3 * time.Second):
		t.Fatal("service did not terminate")
		return nil
	}
}

func TestServeReadinessAndShutdown(t *testing.T) {
	c, _ := configFrom(nil)
	s := newService(c, testLogger())
	addr, cancel, done := startTestService(t, c, s)
	client := &http.Client{Timeout: time.Second, Transport: &http.Transport{DisableKeepAlives: true, Proxy: nil}}
	defer client.CloseIdleConnections()
	for _, path := range []string{"/healthz", "/readyz"} {
		resp, err := client.Get("http://" + addr + path)
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("%s: %d", path, resp.StatusCode)
		}
	}
	cancel()
	if err := waitDone(t, done); err != nil {
		t.Fatal(err)
	}
	if s.ready.Load() {
		t.Fatal("readiness remained true after shutdown")
	}
	if conn, err := net.DialTimeout("tcp", addr, time.Second); err == nil {
		conn.Close()
		t.Fatal("listener still accepts after shutdown")
	}
}

func TestShutdownDrainsOrTimesOut(t *testing.T) {
	for _, timeout := range []bool{false, true} {
		t.Run(map[bool]string{false: "drain", true: "deadline"}[timeout], func(t *testing.T) {
			c, _ := configFrom(nil)
			c.ShutdownTimeout = time.Second
			if timeout {
				c.ShutdownTimeout = 100 * time.Millisecond
			}
			s := newService(c, testLogger())
			entered, release := make(chan struct{}), make(chan struct{})
			defer close(release)
			s.http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				close(entered)
				select {
				case <-release:
				case <-r.Context().Done():
				}
			})
			addr, cancel, done := startTestService(t, c, s)
			clientDone := make(chan error, 1)
			go func() {
				client := &http.Client{Timeout: 2 * time.Second}
				resp, err := client.Get("http://" + addr + "/healthz")
				if err == nil {
					_ = resp.Body.Close()
				}
				client.CloseIdleConnections()
				clientDone <- err
			}()
			select {
			case <-entered:
			case <-time.After(time.Second):
				t.Fatal("request never started")
			}
			cancel()
			if !timeout {
				// Release the in-flight handler after shutdown has begun.
				for deadline := time.Now().Add(time.Second); s.ready.Load(); {
					if time.Now().After(deadline) {
						t.Fatal("readiness did not clear")
					}
					time.Sleep(time.Millisecond)
				}
				select {
				case err := <-done:
					t.Fatalf("shutdown returned before drain: %v", err)
				default:
				}
				release <- struct{}{}
			}
			err := waitDone(t, done)
			if timeout && !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("want deadline, got %v", err)
			}
			if !timeout && err != nil {
				t.Fatal(err)
			}
			select {
			case err := <-clientDone:
				if !timeout && err != nil {
					t.Fatal(err)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("client did not finish")
			}
		})
	}
}

func TestMalformedHTTPAndHeaderLimit(t *testing.T) {
	c, _ := configFrom(nil)
	s := newService(c, testLogger())
	addr, cancel, done := startTestService(t, c, s)
	for _, test := range []struct {
		name, request string
		status        int
	}{
		{"bad-version", "GET /healthz HTTP/invalid\r\nHost: local\r\n\r\n", 400},
		{"bad-header", "GET /healthz HTTP/1.1\r\nHost: local\r\nBad Header: value\r\n\r\n", 400},
		{"ambiguous-length", "GET /healthz HTTP/1.1\r\nHost: local\r\nContent-Length: 1\r\nContent-Length: 2\r\n\r\n", 400},
		{"oversized-header", "GET /healthz HTTP/1.1\r\nHost: local\r\nX-Large: " + strings.Repeat("a", 20<<10) + "\r\n\r\n", 431},
	} {
		t.Run(test.name, func(t *testing.T) {
			conn, err := net.DialTimeout("tcp", addr, time.Second)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			_ = conn.SetDeadline(time.Now().Add(time.Second))
			if _, err := io.WriteString(conn, test.request); err != nil {
				t.Fatal(err)
			}
			resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != test.status {
				t.Fatalf("got %d, want %d", resp.StatusCode, test.status)
			}
		})
	}
	cancel()
	if err := waitDone(t, done); err != nil {
		t.Fatal(err)
	}
}

func TestSlowHeadersAreClosed(t *testing.T) {
	c, _ := configFrom(nil)
	c.ReadHeaderTimeout = 30 * time.Millisecond
	s := newService(c, testLogger())
	addr, cancel, done := startTestService(t, c, s)
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(time.Second))
	if _, err := io.WriteString(conn, "GET /healthz HTTP/1.1\r\nHost:"); err != nil {
		t.Fatal(err)
	}
	_, err = io.ReadAll(conn)
	if err != nil {
		t.Fatalf("server did not close slow headers before client deadline: %v", err)
	}
	cancel()
	if err := waitDone(t, done); err != nil {
		t.Fatal(err)
	}
}

func TestStartupAndServeFailures(t *testing.T) {
	c, _ := configFrom(nil)
	invalid := c
	invalid.Addr = "0.0.0.0:8080"
	if err := Run(context.Background(), invalid, testLogger()); err == nil {
		t.Fatal("invalid config accepted")
	}
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	c.Addr = l.Addr().String()
	if err := Run(context.Background(), c, testLogger()); err == nil {
		t.Fatal("occupied listener accepted")
	}
	_ = l.Close()
	s := newService(c, testLogger())
	if err := s.serve(context.Background(), l, c, testLogger()); err == nil {
		t.Fatal("closed listener did not fail")
	}
	if s.ready.Load() {
		t.Fatal("failed server remains ready")
	}
}
