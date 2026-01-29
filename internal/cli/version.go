package cli

import (
	"fmt"

	"github.com/rishiyaduwanshi/boiler/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Show version information",
	Long: `Display Boiler version information.

Shows current version, build date, and Go version used.`,
	Example: `  # Show version
  bl version

  # Short form
  bl v`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Info())
	},
}
