package list

import (
	"fmt"

	"github.com/rishiyaduwanshi/boiler/internal/config"
	"github.com/rishiyaduwanshi/boiler/internal/ui"
	"github.com/rishiyaduwanshi/boiler/internal/utils"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List snippets or stacks",
	Long: `List all stored snippets and stacks with their version numbers.

By default, shows both snippets and stacks. Use flags to filter by type.
All resources are shown with version numbers included.`,
	Example: `  # List everything
  bl ls

  # List only snippets
  bl ls --snippets

  # List only stacks
  bl ls --stacks`,
	Run: func(cmd *cobra.Command, args []string) {
		st, err := utils.LoadStore(cfg.Paths.Store)
		if err != nil {
			ui.Error.Printf("Error loading store: %v\n", err)
			return
		}

		showAll := !listSnippets && !listStacks

		if listSnippets || showAll {
			ui.Header.Println("\n📄 Snippets:")
			snippets := st.ListSnippets()
			if len(snippets) == 0 {
				ui.Muted.Println("  No snippets found")
			} else {
				for _, name := range snippets {
					ui.Info.Printf("  • %s\n", name)
				}
			}
		}

		if listStacks || showAll {
			ui.Header.Println("\n📦 Stacks:")
			stacks := st.ListStacks()
			if len(stacks) == 0 {
				ui.Muted.Println("  No stacks found")
			} else {
				for _, name := range stacks {
					ui.Info.Printf("  • %s\n", name)
				}
			}
		}

		fmt.Println()
	},
}

var (
	listSnippets bool
	listStacks   bool
)

func init() {
	Cmd.Flags().BoolVarP(&listSnippets, "snippets", "n", false, "List snippets")
	Cmd.Flags().BoolVarP(&listStacks, "stacks", "k", false, "List stacks")
}

var (
	cfg    *config.Config
	logger *utils.Logger
)

func Setup(c *config.Config, l *utils.Logger) {
	cfg = c
	logger = l
}
