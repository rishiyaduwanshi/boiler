package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rishiyaduwanshi/boiler/internal/store"
	"github.com/rishiyaduwanshi/boiler/internal/utils"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [resource]",
	Short: "Add a snippet or stack to current directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]
		logger.Info(fmt.Sprintf("Adding resource: %s", resource))

		if err := addResource(resource); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func addResource(resource string) error {
	st, err := utils.LoadStore(cfg.Paths.Store)
	if err != nil {
		return err
	}

	fullName := utils.ParseResourceName(resource)
	destPath := addTo
	if destPath == "" {
		destPath = "."
	}

	if store.IsSnippet(resource) {
		return addSnippet(st, fullName, destPath)
	}
	return addStack(st, fullName, destPath)
}

func addSnippet(st *store.Store, name, destPath string) error {
	snippetPath, ok := st.GetSnippet(name)
	if !ok {
		return fmt.Errorf(utils.ErrResourceNotFound, "snippet", name)
	}

	if !utils.FileExists(snippetPath) {
		return fmt.Errorf(utils.ErrResourceNotFound, "snippet file", snippetPath)
	}

	destFile := filepath.Join(destPath, filepath.Base(snippetPath))

	if utils.FileExists(destFile) && !addForce {
		return fmt.Errorf(utils.ErrFileAlreadyExists, destFile)
	}

	if err := utils.CopyFile(snippetPath, destFile); err != nil {
		return fmt.Errorf("failed to copy snippet: %w", err)
	}

	fmt.Printf(utils.MsgSnippetAdded, name, destFile)
	logger.Info(fmt.Sprintf("Snippet added: %s -> %s", name, destFile))
	return nil
}

func addStack(st *store.Store, name, destPath string) error {
	stackPath, ok := st.GetStack(name)
	if !ok {
		return fmt.Errorf(utils.ErrResourceNotFound, "stack", name)
	}

	if !utils.IsDirectory(stackPath) {
		return fmt.Errorf(utils.ErrResourceNotFound, "stack directory", stackPath)
	}

	if utils.FileExists(destPath) && destPath != "." && !addForce {
			return fmt.Errorf(utils.ErrDestAlreadyExists, destPath)
	}

	ignorePatterns := []string{"node_modules", ".git", ".DS_Store", "Thumbs.db"}
	if err := utils.CopyDir(stackPath, destPath, ignorePatterns); err != nil {
		return fmt.Errorf("failed to copy stack: %w", err)
	}

	fmt.Printf(utils.MsgStackAdded, name, destPath)
	logger.Info(fmt.Sprintf("Stack added: %s -> %s", name, destPath))
	return nil
}

var (
	addRemote bool
	addTo     string
	addGlobal bool
	addBoth   bool
	addForce  bool
)

func init() {
	addCmd.Flags().BoolVarP(&addRemote, "remote", "r", false, "Fetch from remote registry")
	addCmd.Flags().StringVarP(&addTo, "to", "t", ".", "Destination path")
	addCmd.Flags().BoolVarP(&addGlobal, "global", "g", false, "Add to global store")
	addCmd.Flags().BoolVarP(&addBoth, "both", "b", false, "Add to both local and global")
	addCmd.Flags().BoolVarP(&addForce, FlagForce, FlagForceShort, false, DescForce)
}
