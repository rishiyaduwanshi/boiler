package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rishiyaduwanshi/boiler/internal/ui"
)

func TestNoColorFlagIsInheritedAndParsedForHelp(t *testing.T) {
	root := GetRootCommand()
	old := ui.NoColor
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"ls", "--no-color", "--help"})
	t.Cleanup(func() {
		ui.NoColor = old
		root.PersistentFlags().Lookup("no-color").Changed = false
		root.SetOut(nil)
		root.SetErr(nil)
		root.SetArgs(nil)
	})
	cmd, err := root.ExecuteC()
	if cmd != nil {
		t.Cleanup(func() {
			if help := cmd.Flags().Lookup("help"); help != nil {
				_ = help.Value.Set("false")
				help.Changed = false
			}
		})
	}
	if err != nil {
		t.Fatal(err)
	}
	if !ui.NoColor {
		t.Fatal("--no-color was not applied before help")
	}
	if !strings.Contains(out.String(), "--no-color") {
		t.Fatalf("help does not document the flag: %s", out)
	}
	if got := strings.Count(out.String(), "  bl ls [flags]\n"); got != 1 {
		t.Fatalf("got %d usage lines, want 1: %s", got, out)
	}
	if strings.Contains(out.String(), "\x1b") {
		t.Fatal("help contains ANSI sequences")
	}
}
