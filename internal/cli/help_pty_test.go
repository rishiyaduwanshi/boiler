//go:build !windows

package cli

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"
	"github.com/creack/pty"
	"github.com/rishiyaduwanshi/boiler/internal/ui"
)

func TestHelpStylesUseOutputDestination(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("CLICOLOR", "1")
	t.Setenv("CLICOLOR_FORCE", "1")
	oldNoColor := ui.NoColor
	t.Cleanup(func() { ui.NoColor = oldNoColor })
	for _, tc := range []struct {
		name, noColor string
		flag          bool
	}{
		{name: "color"},
		{name: "flag", flag: true},
		{name: "environment", noColor: "1"},
		{name: "environment zero", noColor: "0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("NO_COLOR", tc.noColor)
			ui.NoColor = tc.flag
			master, slave, err := pty.Open()
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = master.Close(); _ = slave.Close() })
			state, err := term.MakeRaw(slave.Fd())
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = term.Restore(slave.Fd(), state) })
			type result struct {
				text string
				err  error
			}
			output := make(chan result, 1)
			go func() {
				text, err := bufio.NewReader(master).ReadString('\x00')
				output <- result{text, err}
			}()
			root, list := helpFixture(false)
			configureHelp(root)
			// stdout is untouched: rendering must honor Cobra's actual writer.
			root.SetOut(slave)
			if err := list.Help(); err != nil {
				t.Fatal(err)
			}
			if _, err := slave.Write([]byte{0}); err != nil {
				t.Fatal(err)
			}
			select {
			case got := <-output:
				if got.err != nil {
					t.Fatal(got.err)
				}
				colored := !tc.flag && tc.noColor == ""
				if strings.Contains(got.text, "\x1b[") != colored {
					t.Fatalf("want colored=%v, got %q", colored, got.text)
				}
				for _, section := range []string{"Usage:", "Aliases:", "Examples:", "Flags:", "Global Flags:"} {
					if !strings.Contains(got.text, ui.Header.RenderTo(slave, section)+"\n") {
						t.Fatalf("section %s does not use ui.Header: %q", section, got.text)
					}
				}
				if strings.Count(ansi.Strip(got.text), "  bl ls [flags]\n") != 1 {
					t.Fatalf("expected one usage line: %q", got.text)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("timed out reading help from PTY")
			}
		})
	}
}
