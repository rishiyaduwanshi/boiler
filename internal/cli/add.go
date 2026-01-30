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
	Long: `Add a stored snippet or stack to your current directory.

The command copies resources from your store. For snippets with a single version,
you can use just the name (e.g., 'errorHandler' will auto-select version 1).
For multiple versions, you'll be prompted to choose.

Template Variables:
  Snippets can contain template variables using the format: bl__VAR_NAME
  When adding a snippet with variables, you'll be prompted to provide values:
    - Default values are shown in brackets (from __var declarations)
    - Press Enter to use default or type a custom value
    - Variables are replaced and metadata comments are removed in the final file

Stacks are also versioned and can be added by name or with explicit version.`,
	Example: `  # Add snippet (auto-detects if single version)
  bl add errorHandler

  # Add snippet with template variables
  bl add apiClient
  # Prompts: bl__API_URL [http://localhost:3000]: https://api.example.com
  #          bl__API_KEY [your-key]: abc123xyz
  # Output: Clean file with variables replaced, no metadata comments

  # Add specific version
  bl add logger@2.js

  # Add to specific directory
  bl add config --to ./src/utils

  # Add stack
  bl add express-api@1

  # Force overwrite
  bl add middleware --force`,
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

	destPath := addTo
	if destPath == "" {
		destPath = "."
	}

	// Parse resource name to extract parts
	baseName, version, ext := store.ParseResourceName(resource)

	// If it's a snippet (has extension)
	if ext != "" {
		// If version is explicitly provided, use it directly
		if version != "" {
			fullName := baseName + "@" + version + ext
			return addSnippet(st, fullName, destPath)
		}

		// No version specified - find matching snippets by name and extension
		matchingSnippets := findMatchingSnippetsByNameAndExt(st, baseName, ext)
		
		if len(matchingSnippets) == 0 {
			return fmt.Errorf(utils.ErrResourceNotFound, "snippet", resource)
		}

		// If only one version exists, use it automatically
		if len(matchingSnippets) == 1 {
			return addSnippet(st, matchingSnippets[0], destPath)
		}

		// Multiple versions - prompt user to choose
		fmt.Printf("Multiple versions found for '%s%s':\n", baseName, ext)
		for i, name := range matchingSnippets {
			fmt.Printf("  %d. %s\n", i+1, name)
		}

		choice, err := utils.Prompt(fmt.Sprintf("Enter version number (1-%d): ", len(matchingSnippets)))
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		var selectedIdx int
		fmt.Sscanf(choice, "%d", &selectedIdx)
		if selectedIdx < 1 || selectedIdx > len(matchingSnippets) {
			return fmt.Errorf("invalid choice")
		}

		return addSnippet(st, matchingSnippets[selectedIdx-1], destPath)
	}

	// No extension - could be stack or snippet name without version/extension
	// First check if it exists as a stack
	_, stackExists := st.GetStack(resource)
	if stackExists {
		return addStack(st, resource, destPath)
	}

	// Not a stack, try to find matching snippets by base name only
	matchingSnippets := findMatchingSnippets(st, baseName)
	if len(matchingSnippets) == 0 {
		return fmt.Errorf(utils.ErrResourceNotFound, "stack or snippet", resource)
	}

	// If only one version exists, use it automatically
	if len(matchingSnippets) == 1 {
		return addSnippet(st, matchingSnippets[0], destPath)
	}

	// Multiple versions - prompt user to choose
	fmt.Printf("Multiple versions found for '%s':\n", baseName)
	for i, name := range matchingSnippets {
		fmt.Printf("  %d. %s\n", i+1, name)
	}

	choice, err := utils.Prompt(fmt.Sprintf("Enter version number (1-%d): ", len(matchingSnippets)))
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Parse choice
	var selectedIdx int
	fmt.Sscanf(choice, "%d", &selectedIdx)
	if selectedIdx < 1 || selectedIdx > len(matchingSnippets) {
		return fmt.Errorf("invalid choice")
	}

	return addSnippet(st, matchingSnippets[selectedIdx-1], destPath)
}

// findMatchingSnippets finds all snippets that match the given name (without version/extension)
func findMatchingSnippets(st *store.Store, name string) []string {
	allSnippets := st.ListSnippets()
	var matches []string

	for _, snippet := range allSnippets {
		snippetName, _, _ := store.ParseResourceName(snippet)
		if snippetName == name {
			matches = append(matches, snippet)
		}
	}

	return matches
}

// findMatchingSnippetsByNameAndExt finds all snippets matching both name and extension
func findMatchingSnippetsByNameAndExt(st *store.Store, name, ext string) []string {
	allSnippets := st.ListSnippets()
	var matches []string

	for _, snippet := range allSnippets {
		snippetName, _, snippetExt := store.ParseResourceName(snippet)
		if snippetName == name && snippetExt == ext {
			matches = append(matches, snippet)
		}
	}

	return matches
}

func addSnippet(st *store.Store, name, destPath string) error {
	snippetPath, ok := st.GetSnippet(name)
	if !ok {
		return fmt.Errorf(utils.ErrResourceNotFound, "snippet", name)
	}

	if !utils.FileExists(snippetPath) {
		return fmt.Errorf(utils.ErrResourceNotFound, "snippet file", snippetPath)
	}

	// Parse snippet metadata to check for variables
	meta, err := utils.ParseSnippetMetadata(snippetPath)
	if err != nil {
		return fmt.Errorf("failed to parse snippet metadata: %w", err)
	}

	// Prompt user for variable values if variables exist
	varReplacements := make(map[string]string)
	if len(meta.Variables) > 0 {
		fmt.Println("Template variables found:")
		for varName, defaultValue := range meta.Variables {
			prompt := fmt.Sprintf("  %s", varName)
			value, err := utils.PromptWithDefault(prompt, defaultValue)
			if err != nil {
				return fmt.Errorf("failed to read variable input: %w", err)
			}
			varReplacements[varName] = value
		}
	}

	// Extract base name without version: errorHandler@1.js -> errorHandler.js
	baseName, _, ext := store.ParseResourceName(name)
	destFileName := baseName + ext
	destFile := filepath.Join(destPath, destFileName)

	if utils.FileExists(destFile) && !addForce {
		return fmt.Errorf(utils.ErrFileAlreadyExists, destFile)
	}

	// Copy file with variable replacement
	if err := utils.CopyFileWithVariables(snippetPath, destFile, varReplacements); err != nil {
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
