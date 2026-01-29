package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/rishiyaduwanshi/boiler/internal/utils"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for snippets or stacks",
	Long: `Search for resources in your store by name.

Searches both snippets and stacks by default. Use flags to filter:
  - Use -s or --snippets to search only snippets
  - Use -k or --stacks to search only stacks

Search is case-insensitive and matches partial names.`,
	Example: `  # Search for anything with 'error'
  bl search error

  # Search only snippets
  bl search logger --snippets

  # Search only stacks
  bl search express --stacks`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		logger.Info(fmt.Sprintf("Searching for: %s", query))

		if err := searchResources(query); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func searchResources(query string) error {
	st, err := utils.LoadStore(cfg.Paths.Store)
	if err != nil {
		return err
	}

	query = strings.ToLower(query)
	foundAny := false

	// Search snippets
	if !searchStacks {
		snippets := st.ListSnippets()
		matches := []string{}
		for _, name := range snippets {
			if strings.Contains(strings.ToLower(name), query) {
				matches = append(matches, name)
			}
		}

		if len(matches) > 0 {
			foundAny = true
			fmt.Println("\n📄 Snippets:")
			for _, name := range matches {
				fmt.Printf("  • %s\n", name)
			}
		}
	}

	// Search stacks
	if !searchSnippets {
		stacks := st.ListStacks()
		matches := []string{}
		for _, name := range stacks {
			if strings.Contains(strings.ToLower(name), query) {
				matches = append(matches, name)
			}
		}

		if len(matches) > 0 {
			foundAny = true
			fmt.Println("\n📦 Stacks:")
			for _, name := range matches {
				fmt.Printf("  • %s\n", name)
			}
		}
	}

	if !foundAny {
		fmt.Printf("No results found for '%s'\n", query)
	} else {
		fmt.Println()
	}

	return nil
}

var (
	searchSnippets bool
	searchStacks   bool
	searchRemote   bool
)

func init() {
	searchCmd.Flags().BoolVarP(&searchSnippets, "snippets", "n", false, "Search only snippets")
	searchCmd.Flags().BoolVarP(&searchStacks, "stacks", "k", false, "Search only stacks")
	searchCmd.Flags().BoolVarP(&searchRemote, "remote", "r", false, "Search remote registry")
}
