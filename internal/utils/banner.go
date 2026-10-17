package utils

import (
	"fmt"

	"github.com/rishiyaduwanshi/boiler/internal/ui"
	"github.com/rishiyaduwanshi/boiler/pkg/version"
)

func ShowBanner() {
	art := `
██████╗  ██████╗ ██╗██╗     ███████╗██████╗ 
██╔══██╗██╔═══██╗██║██║     ██╔════╝██╔══██╗
██████╔╝██║   ██║██║██║     █████╗  ██████╔╝
██╔══██╗██║   ██║██║██║     ██╔══╝  ██╔══██╗
██████╔╝╚██████╔╝██║███████╗███████╗██║  ██║
╚═════╝  ╚═════╝ ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝                                   
`

	ui.Info.Println(art)
	ui.Muted.Println("Code Snippet & Stack Manager")
	ui.Info.Printf("Version: ")
	ui.Warning.Println(version.Version)
	ui.Muted.Println("by Abhinav Prakash")
}

func ShowQuickHelp() {
	ui.Header.Println("\nQuick Commands:")
	ui.Info.Println("  bl add <resource>        Add snippet or stack")
	ui.Info.Println("  bl store <path>          Store as snippet/stack")
	ui.Info.Println("  bl ls                   List all resources")
	ui.Info.Println("  bl --help                Show full help")
	fmt.Println()
}
