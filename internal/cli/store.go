package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rishiyaduwanshi/boiler/internal/store"
	"github.com/rishiyaduwanshi/boiler/internal/utils"
	"github.com/spf13/cobra"
)

var (
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
)

var storeCmd = &cobra.Command{
	Use:   "store [path]",
	Short: "Store a folder/file as snippet or stack",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		logger.Info(fmt.Sprintf("Storing: %s", path))

		if err := storeResource(path); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func storeResource(path string) error {
	// Check if path exists
	if !utils.FileExists(path) {
		return fmt.Errorf("path '%s' does not exist", path)
	}

	// Auto-detect type if not provided
	isDir := utils.IsDirectory(path)
	autoDetectedType := "snippet"
	if isDir {
		autoDetectedType = "stack"
	}

	// Use flags to override auto-detection
	if storeAsSnippet {
		autoDetectedType = "snippet"
	} else if storeAsStack {
		autoDetectedType = "stack"
	}

	// Get base name if not provided
	if storeName == "" {
		storeName = filepath.Base(path)
		// Remove extension for snippets to get clean name
		if !isDir {
			storeName = strings.TrimSuffix(storeName, filepath.Ext(storeName))
		}
	}

	st := store.NewStore(cfg.Paths.Store)
	if err := st.Load(); err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	if autoDetectedType == "snippet" {
		return storeSnippet(st, path)
	}
	return storeStack(st, path)
}

func storeSnippet(st *store.Store, path string) error {
	// Must be a file
	if utils.IsDirectory(path) {
		return fmt.Errorf("snippet must be a file, not a directory")
	}

	// Get extension
	ext := filepath.Ext(path)
	if ext == "" {
		return fmt.Errorf("snippet file must have an extension")
	}

	// Determine the language directory based on extension
	langDir := strings.TrimPrefix(ext, ".")
	snippetDir := filepath.Join(cfg.Paths.Snippets, langDir)
	if err := utils.EnsureDir(snippetDir); err != nil {
		return fmt.Errorf("failed to create snippet directory: %w", err)
	}

	// Check if any version exists in metadata
	existingVersions := findExistingVersions(st, storeName, ext, true)
	var version int
	var fullName string
	var destPath string

	if len(existingVersions) > 0 {
		// Already exists - show warning and ask user
		latestVersion := existingVersions[len(existingVersions)-1]
		fmt.Println(warningStyle.Render(fmt.Sprintf("⚠ Warning: snippet '%s' already exists in store", storeName)))
		choice, err := utils.Prompt("Do you want to create a new version (v) or overwrite it (o)? ")
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		choice = strings.ToLower(strings.TrimSpace(choice))
		if choice == "v" {
			// Create new version
			version = latestVersion + 1
			fullName = fmt.Sprintf("%s@%d%s", storeName, version, ext)
			destPath = filepath.Join(snippetDir, filepath.Base(fullName))
		} else if choice == "o" {
			// Overwrite latest version
			version = latestVersion
			fullName = fmt.Sprintf("%s@%d%s", storeName, version, ext)
			destPath = filepath.Join(snippetDir, filepath.Base(fullName))
			// Remove old entry from metadata
			if err := st.RemoveSnippet(fullName); err != nil {
				return fmt.Errorf("failed to remove old snippet: %w", err)
			}
			// Remove old file
			if utils.FileExists(destPath) {
				if err := os.Remove(destPath); err != nil {
					return fmt.Errorf("failed to remove old file: %w", err)
				}
			}
		} else {
			return fmt.Errorf("cancelled")
		}
	} else {
		// First version
		version = 1
		fullName = fmt.Sprintf("%s@%d%s", storeName, version, ext)
		destPath = filepath.Join(snippetDir, filepath.Base(fullName))
	}

	// Copy file
	if err := utils.CopyFile(path, destPath); err != nil {
		return fmt.Errorf("failed to copy snippet: %w", err)
	}

	// Add to metadata
	if err := st.AddSnippet(fullName, destPath); err != nil {
		return fmt.Errorf("failed to update metadata: %w", err)
	}

	fmt.Printf("✓ Stored snippet '%s' at %s\n", fullName, destPath)
	logger.Info(fmt.Sprintf("Snippet stored: %s -> %s", path, destPath))
	return nil
}

func storeStack(st *store.Store, path string) error {
	// Must be a directory
	if !utils.IsDirectory(path) {
		return fmt.Errorf("stack must be a directory, not a file")
	}

	// Check if any version exists in metadata
	existingVersions := findExistingVersions(st, storeName, "", false)
	var version int
	var fullName string
	var stackDir string

	if len(existingVersions) > 0 {
		// Already exists - show warning and ask user
		latestVersion := existingVersions[len(existingVersions)-1]
		fmt.Println(warningStyle.Render(fmt.Sprintf("⚠ Warning: stack '%s' already exists in store", storeName)))
		choice, err := utils.Prompt("Do you want to create a new version (v) or overwrite it (o)? ")
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		choice = strings.ToLower(strings.TrimSpace(choice))
		if choice == "v" {
			// Create new version
			version = latestVersion + 1
			fullName = fmt.Sprintf("%s@%d", storeName, version)
			stackDir = filepath.Join(cfg.Paths.Stacks, fullName)
		} else if choice == "o" {
			// Overwrite latest version
			version = latestVersion
			fullName = fmt.Sprintf("%s@%d", storeName, version)
			stackDir = filepath.Join(cfg.Paths.Stacks, fullName)
			// Remove old entry from metadata
			if err := st.RemoveStack(fullName); err != nil {
				return fmt.Errorf("failed to remove old stack: %w", err)
			}
			// Remove old directory
			if utils.IsDirectory(stackDir) {
				if err := os.RemoveAll(stackDir); err != nil {
					return fmt.Errorf("failed to remove old directory: %w", err)
				}
			}
		} else {
			return fmt.Errorf("cancelled")
		}
	} else {
		// First version
		version = 1
		fullName = fmt.Sprintf("%s@%d", storeName, version)
		stackDir = filepath.Join(cfg.Paths.Stacks, fullName)
	}

	// Copy directory
	ignorePatterns := []string{"node_modules", ".git", ".DS_Store", "Thumbs.db", "vendor", "__pycache__", ".vscode"}
	if err := utils.CopyDir(path, stackDir, ignorePatterns); err != nil {
		return fmt.Errorf("failed to copy stack: %w", err)
	}

	// Add to metadata
	if err := st.AddStack(fullName, stackDir); err != nil {
		return fmt.Errorf("failed to update metadata: %w", err)
	}

	fmt.Printf("✓ Stored stack '%s' at %s\n", fullName, stackDir)
	logger.Info(fmt.Sprintf("Stack stored: %s -> %s", path, stackDir))
	return nil
}

// findExistingVersions returns all version numbers for a given base name
func findExistingVersions(st *store.Store, baseName string, ext string, isSnippet bool) []int {
	versions := []int{}
	pattern := regexp.MustCompile(fmt.Sprintf(`^%s@(\d+)`, regexp.QuoteMeta(baseName)))

	var items []string
	if isSnippet {
		items = st.ListSnippets()
	} else {
		items = st.ListStacks()
	}

	for _, item := range items {
		// Remove extension if snippet
		itemName := item
		if isSnippet && ext != "" {
			itemName = strings.TrimSuffix(item, ext)
		}

		matches := pattern.FindStringSubmatch(itemName)
		if len(matches) > 1 {
			if ver, err := strconv.Atoi(matches[1]); err == nil {
				versions = append(versions, ver)
			}
		}
	}

	return versions
}

var (
	storeName        string
	storeAsSnippet   bool
	storeAsStack     bool
	storeDescription string
)

func init() {
	storeCmd.Flags().StringVarP(&storeName, "name", "", "", "Name for the resource (auto-detected from path if not provided)")
	storeCmd.Flags().BoolVarP(&storeAsSnippet, "snippet", "n", false, "Force store as snippet")
	storeCmd.Flags().BoolVarP(&storeAsStack, "stack", "k", false, "Force store as stack")
	storeCmd.Flags().StringVarP(&storeDescription, "description", "d", "", "Description")
}
