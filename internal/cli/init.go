package cli

import (
	"fmt"
	"os"

	"github.com/rishiyaduwanshi/boiler/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize boiler in current directory",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Initializing boiler")

		if err := initBoiler(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func initBoiler() error {
	// Initialize configuration
	if err := cfg.InitializeDirs(); err != nil {
		return fmt.Errorf("failed to initialize directories: %w", err)
	}

	// Create config file if it doesn't exist
	configPath, err := config.ConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Save current config
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("✓ Boiler initialized successfully!\n\n")
	fmt.Printf("Configuration: %s\n", configPath)
	fmt.Printf("Store path:    %s\n", cfg.Paths.Store)
	fmt.Printf("Snippets:      %s\n", cfg.Paths.Snippets)
	fmt.Printf("Stacks:        %s\n", cfg.Paths.Stacks)
	fmt.Printf("\nUse 'bl store' to add snippets/stacks or 'bl add' to use them.\n")

	logger.Info("Boiler initialization complete")
	return nil
}

var (
	initGlobal bool
)

func init() {
	initCmd.Flags().BoolVarP(&initGlobal, "global", "g", false, "Initialize global configuration")
}
