package path

import (
	"github.com/rishiyaduwanshi/boiler/internal/config"
	"github.com/rishiyaduwanshi/boiler/internal/ui"
	"github.com/rishiyaduwanshi/boiler/internal/utils"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "path",
	Short: "Show boiler installation path",
	Long: `Display all Boiler installation paths.

Shows:
  - Root - Main Boiler directory
  - Store - Where resources are stored
  - Snippets - Snippet storage location
  - Stacks - Stack storage location
  - Logs - Log file directory
  - Bin - Executable location`,
	Example: `  # Show all paths
  bl path

  # Use in scripts
  cd $(bl path | grep Store | cut -d: -f2)`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.Info.Printf("Boiler root: %s\n", ui.Path.Render(cfg.Paths.Root))
		ui.Info.Printf("Store:       %s\n", ui.Path.Render(cfg.Paths.Store))
		ui.Info.Printf("Snippets:    %s\n", ui.Path.Render(cfg.Paths.Snippets))
		ui.Info.Printf("Stacks:      %s\n", ui.Path.Render(cfg.Paths.Stacks))
		ui.Info.Printf("Logs:        %s\n", ui.Path.Render(cfg.Paths.Logs))
		ui.Info.Printf("Bin:         %s\n", ui.Path.Render(cfg.Paths.Bin))
	},
}

var (
	cfg    *config.Config
	logger *utils.Logger
)

func Setup(c *config.Config, l *utils.Logger) {
	cfg = c
	logger = l
}
