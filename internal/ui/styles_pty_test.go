//go:build !windows

package ui

import (
	"bufio"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"
	"github.com/creack/pty"
)

// TestStylesRenderColorOnTTY verifies the colored rendering branch by driving a
// real pseudo-terminal. It is gated to non-Windows platforms because
// github.com/creack/pty does not support opening PTYs on Windows.
func TestStylesRenderColorOnTTY(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "")
	t.Setenv("CLICOLOR", "1")
	t.Setenv("CLICOLOR_FORCE", "1")
	oldNoColor := NoColor
	NoColor = false
	t.Cleanup(func() { NoColor = oldNoColor })

	master, slave, err := pty.Open()
	if err != nil {
		t.Fatalf("open PTY: %v", err)
	}
	t.Cleanup(func() { _ = master.Close() })
	t.Cleanup(func() { _ = slave.Close() })

	// Prevent the terminal driver from translating LF to CRLF, so the test
	// checks the exact newlines written by the renderer.
	state, err := term.MakeRaw(slave.Fd())
	if err != nil {
		t.Fatalf("make PTY raw: %v", err)
	}
	t.Cleanup(func() { _ = term.Restore(slave.Fd(), state) })
	if !colorEnabled(slave) {
		t.Fatal("color must be enabled for the PTY slave")
	}

	// Render uses stdout rather than the writer passed to Fprintf.
	oldStdout := os.Stdout
	os.Stdout = slave
	t.Cleanup(func() { os.Stdout = oldStdout })
	path := Path.Render("example.go")
	if !strings.Contains(path, "\x1b[") {
		t.Fatalf("rendered path has no ANSI sequences: %q", path)
	}
	if got := ansi.Strip(path); got != "example.go" {
		t.Fatalf("plain rendered path = %q, want %q", got, "example.go")
	}

	type readResult struct {
		text string
		err  error
	}
	output := make(chan readResult, 1)
	go func() {
		text, err := bufio.NewReader(master).ReadString('\x00')
		output <- readResult{text: text, err: err}
	}()

	Success.Fprintf(slave, "\n\n  saved %s\n\n", path)
	// Delimit the output without relying on platform-specific PTY EOF behavior.
	if _, err := slave.Write([]byte{0}); err != nil {
		t.Fatalf("write output delimiter: %v", err)
	}
	select {
	case result := <-output:
		if result.err != nil {
			t.Fatalf("read PTY output: %v", result.err)
		}
		got := strings.TrimSuffix(result.text, "\x00")
		if !strings.Contains(got, "\x1b[") {
			t.Fatalf("colored output has no ANSI sequences: %q", got)
		}
		if !strings.HasPrefix(got, "\n\n") || !strings.HasSuffix(got, "\n\n") {
			t.Fatalf("leading/trailing newlines were not preserved: %q", got)
		}
		if plain, want := ansi.Strip(got), "\n\n  saved example.go\n\n"; plain != want {
			t.Fatalf("plain nested output = %q, want %q (colored: %q)", plain, want, got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out reading PTY output")
	}
}
