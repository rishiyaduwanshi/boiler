package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/rishiyaduwanshi/boiler/internal/store"
	"github.com/rishiyaduwanshi/boiler/internal/utils"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean [resource]",
	Short: "Clean snippets, stacks, or store",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			resource := args[0]
			logger.Info(fmt.Sprintf("Cleaning resource: %s", resource))
			if err := cleanResource(resource); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		if cleanAll {
			if err := cleanAllResources(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		if cleanSnippets {
			if err := cleanAllSnippets(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		if cleanStacks {
			if err := cleanAllStacks(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		if err := interactiveClean(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var (
	cleanAll      bool
	cleanSnippets bool
	cleanStacks   bool
)

func cleanResource(resource string) error {
	st, err := utils.LoadStore(cfg.Paths.Store)
	if err != nil {
		return err
	}

	fullName := utils.ParseResourceName(resource)

	if store.IsSnippet(resource) {
		return cleanSnippet(st, fullName)
	}
	return cleanStack(st, fullName)
}

func cleanSnippet(st *store.Store, name string) error {
	path, ok := st.GetSnippet(name)
	if !ok {
		return fmt.Errorf(utils.ErrResourceNotFound, "snippet", name)
	}

	if !utils.ConfirmAction(fmt.Sprintf(utils.MsgPromptConfirmRemove, "snippet", name)) {
		fmt.Println(utils.MsgCancelled)
		return nil
	}

	if utils.FileExists(path) {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("failed to remove file: %w", err)
		}
	}

	if err := st.RemoveSnippet(name); err != nil {
		return fmt.Errorf("failed to update metadata: %w", err)
	}

	fmt.Printf(utils.MsgSnippetRemoved, name)
	logger.Info(fmt.Sprintf("Snippet removed: %s", name))
	return nil
}

func cleanStack(st *store.Store, name string) error {
	path, ok := st.GetStack(name)
	if !ok {
		return fmt.Errorf(utils.ErrResourceNotFound, "stack", name)
	}

	if !utils.ConfirmAction(fmt.Sprintf(utils.MsgPromptConfirmRemove, "stack", name)) {
		fmt.Println(utils.MsgCancelled)
		return nil
	}

	if utils.IsDirectory(path) {
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove directory: %w", err)
		}
	}

	if err := st.RemoveStack(name); err != nil {
		return fmt.Errorf("failed to update metadata: %w", err)
	}

	fmt.Printf(utils.MsgStackRemoved, name)
	logger.Info(fmt.Sprintf("Stack removed: %s", name))
	return nil
}

func cleanAllResources() error {
	st, err := utils.LoadStore(cfg.Paths.Store)
	if err != nil {
		return err
	}

	fmt.Println(utils.MsgPromptConfirmCleanAll)
	if !utils.ConfirmAction("Are you sure? (y/N): ") {
		fmt.Println(utils.MsgCancelled)
		return nil
	}

	snippets := st.ListSnippets()
	for _, name := range snippets {
		path, _ := st.GetSnippet(name)
		if utils.FileExists(path) {
			os.Remove(path)
		}
		st.RemoveSnippet(name)
	}

	stacks := st.ListStacks()
	for _, name := range stacks {
		path, _ := st.GetStack(name)
		if utils.IsDirectory(path) {
			os.RemoveAll(path)
		}
		st.RemoveStack(name)
	}

	fmt.Printf("✓ Removed %d snippets and %d stacks\n", len(snippets), len(stacks))
	logger.Info("All resources cleaned")
	return nil
}

func cleanAllSnippets() error {
	st, err := utils.LoadStore(cfg.Paths.Store)
	if err != nil {
		return err
	}

	snippets := st.ListSnippets()
	if len(snippets) == 0 {
		fmt.Println(utils.MsgNoSnippets)
		return nil
	}

	if !utils.ConfirmAction(fmt.Sprintf("Remove %d snippets? (y/N): ", len(snippets))) {
		fmt.Println(utils.MsgCancelled)
		return nil
	}

	for _, name := range snippets {
		path, _ := st.GetSnippet(name)
		if utils.FileExists(path) {
			os.Remove(path)
		}
		st.RemoveSnippet(name)
	}

	fmt.Printf("✓ Removed %d snippets\n", len(snippets))
	logger.Info(fmt.Sprintf("Cleaned %d snippets", len(snippets)))
	return nil
}

func cleanAllStacks() error {
	st, err := utils.LoadStore(cfg.Paths.Store)
	if err != nil {
		return err
	}

	stacks := st.ListStacks()
	if len(stacks) == 0 {
		fmt.Println(utils.MsgNoStacks)
		return nil
	}

	if !utils.ConfirmAction(fmt.Sprintf("Remove %d stacks? (y/N): ", len(stacks))) {
		fmt.Println(utils.MsgCancelled)
		return nil
	}

	for _, name := range stacks {
		path, _ := st.GetStack(name)
		if utils.IsDirectory(path) {
			os.RemoveAll(path)
		}
		st.RemoveStack(name)
	}

	fmt.Printf("✓ Removed %d stacks\n", len(stacks))
	logger.Info(fmt.Sprintf("Cleaned %d stacks", len(stacks)))
	return nil
}

func interactiveClean() error {
	fmt.Println("\nSelect action:")
	fmt.Println("  k - clean all stacks")
	fmt.Println("  n - clean all snippets")
	fmt.Println("  a - clean all resources")
	fmt.Println("  q - quit")
	fmt.Print("\nChoice: ")

	var choice string
	fmt.Scanln(&choice)
	choice = strings.ToLower(strings.TrimSpace(choice))

	switch choice {
	case "k":
		return cleanAllStacks()
	case "n":
		return cleanAllSnippets()
	case "a":
		return cleanAllResources()
	case "q":
		fmt.Println(utils.MsgCancelled)
		return nil
	default:
		fmt.Println("Invalid choice")
		return nil
	}
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVarP(&cleanAll, FlagAll, FlagAllShort, false, DescCleanAll)
	cleanCmd.Flags().BoolVarP(&cleanSnippets, FlagSnippets, FlagSnippetsShort, false, DescSnippetsOnly)
	cleanCmd.Flags().BoolVarP(&cleanStacks, FlagStacks, FlagStacksShort, false, DescStacksOnly)
}
