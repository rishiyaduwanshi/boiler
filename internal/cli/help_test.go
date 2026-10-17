package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func helpFixture(grouped bool) (*cobra.Command, *cobra.Command) {
	root := &cobra.Command{Use: "bl", Short: "Boiler", Run: func(*cobra.Command, []string) {}}
	root.PersistentFlags().Bool("no-color", false, "Disable colors")
	list := &cobra.Command{
		Use: "ls", Aliases: []string{"list", "items", "resources", "all"},
		Short: "List resources", Long: "List stored snippets and stacks.",
		Example: "  # List everything\n  bl ls", Run: func(*cobra.Command, []string) {},
	}
	list.Flags().BoolP("snippets", "n", false, "List snippets")
	list.InitDefaultHelpFlag()
	root.AddCommand(list,
		&cobra.Command{Use: "hidden", Hidden: true, Run: func(*cobra.Command, []string) {}},
		&cobra.Command{Use: "old", Deprecated: "use ls", Run: func(*cobra.Command, []string) {}},
		&cobra.Command{Use: "topic", Short: "A help topic"},
		&cobra.Command{Use: "other", Short: "Another command", Run: func(*cobra.Command, []string) {}},
	)
	if grouped {
		root.AddGroup(&cobra.Group{ID: "resources", Title: "Resource Commands:"})
		list.GroupID = "resources"
	}
	root.InitDefaultHelpCmd()
	root.InitDefaultHelpFlag()
	return root, list
}

func TestHelpPreservesCobraPlainLayout(t *testing.T) {
	for _, grouped := range []bool{false, true} {
		for _, name := range []string{"bl", "ls"} {
			root, list := helpFixture(grouped)
			cmd := root
			if name == "ls" {
				cmd = list
			}
			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&out)
			// Capture Cobra's layout before installing Boiler's renderer.
			if err := cmd.Help(); err != nil {
				t.Fatal(err)
			}
			wantHelp := out.String()
			wantUsage := cmd.UsageString()
			configureHelp(cmd)
			for i := 0; i < 2; i++ {
				out.Reset()
				if err := cmd.Help(); err != nil {
					t.Fatal(err)
				}
				if got := out.String(); got != wantHelp {
					t.Fatalf("%s grouped=%v help mismatch:\ngot:\n%s\nwant:\n%s", cmd.Name(), grouped, got, wantHelp)
				}
				if got := cmd.UsageString(); got != wantUsage {
					t.Fatalf("%s usage mismatch:\ngot:\n%s\nwant:\n%s", cmd.Name(), got, wantUsage)
				}
			}
		}
	}
}

func TestHelpUsageIsNotRepeatedForAliases(t *testing.T) {
	for _, name := range []string{"ls", "list", "items", "resources", "all"} {
		t.Run(name, func(t *testing.T) {
			root, _ := helpFixture(false)
			configureHelp(root)
			var out bytes.Buffer
			root.SetOut(&out)
			root.SetErr(&out)
			root.SetArgs([]string{name, "--help"})
			if err := root.Execute(); err != nil {
				t.Fatal(err)
			}
			if got := strings.Count(out.String(), "  bl ls [flags]\n"); got != 1 {
				t.Fatalf("got %d usage lines, want 1:\n%s", got, &out)
			}
			for _, section := range []string{"Usage:", "Aliases:", "Examples:", "Flags:", "Global Flags:"} {
				if strings.Count(out.String(), "\n"+section+"\n") != 1 {
					t.Fatalf("missing or repeated %s section:\n%s", section, &out)
				}
			}
		})
	}
}

type failingHelpWriter struct{ err error }

func (w failingHelpWriter) Write([]byte) (int, error) { return 0, w.err }

func TestUsageReturnsWriterError(t *testing.T) {
	root, _ := helpFixture(false)
	configureHelp(root)
	want := errors.New("write failed")
	root.SetOut(failingHelpWriter{want})
	if err := root.Usage(); !errors.Is(err, want) {
		t.Fatalf("Usage error = %v, want %v", err, want)
	}
}
