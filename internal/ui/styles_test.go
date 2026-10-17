package ui

import (
	"bytes"
	"os"
	"testing"
)

func TestStylesPreservePlainOutput(t *testing.T) {
	styles := map[string]Style{
		"error": Error, "success": Success, "info": Info, "warning": Warning,
		"path": Path, "header": Header, "muted": Muted,
	}
	for name, style := range styles {
		t.Run(name, func(t *testing.T) {
			var out bytes.Buffer
			style.Fprintf(&out, "\n  %s: %d%%\n\n", "hello", 100)
			if got, want := out.String(), "\n  hello: 100%\n\n"; got != want {
				t.Fatalf("output = %q, want %q", got, want)
			}
		})
	}
}

func TestColorOptOut(t *testing.T) {
	old := NoColor
	t.Cleanup(func() { NoColor = old })
	for _, tc := range []struct {
		name, noColor, term string
		flag                bool
	}{
		{name: "flag", term: "xterm-256color", flag: true},
		{name: "environment", noColor: "1", term: "xterm-256color"},
		{name: "environment zero is nonempty", noColor: "0", term: "xterm-256color"},
		{name: "dumb terminal", term: "dumb"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			NoColor = tc.flag
			t.Setenv("NO_COLOR", tc.noColor)
			t.Setenv("TERM", tc.term)
			t.Setenv("CLICOLOR_FORCE", "1")
			if colorEnabled(os.Stdout) {
				t.Fatal("color enabled despite explicit opt-out")
			}
			if got := Header.Render("Title"); got != "Title" {
				t.Fatalf("disabled header = %q", got)
			}
		})
	}
}

func TestRedirectedOutputIgnoresForcedColor(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("CLICOLOR_FORCE", "1")
	file, err := os.CreateTemp(t.TempDir(), "output")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if colorEnabled(file) {
		t.Fatal("regular files must not receive ANSI sequences")
	}
	Success.Fprintf(file, "saved %s\n", "example.go")
	if _, err := file.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if _, err := out.ReadFrom(file); err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), "saved example.go\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}
