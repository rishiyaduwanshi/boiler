package cli

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/rishiyaduwanshi/boiler/internal/ui"
	"github.com/spf13/cobra"
)

// Keep Cobra's section order and flag alignment, but render semantic styles
// through internal/ui. UseLine is rendered once, outside all command/alias loops.
const usageTemplate = `{{heading "Usage:"}}{{if .Runnable}}
  {{info .UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{info (printf "%s [command]" .CommandPath)}}{{end}}{{if .Aliases}}

{{heading "Aliases:"}}
  {{info .NameAndAliases}}{{end}}{{if .HasExample}}

{{heading "Examples:"}}
{{examples .Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if not .Groups}}

{{heading "Available Commands:"}}{{range $cmds}}{{if or .IsAvailableCommand (eq .Name "help")}}
  {{command .Name .NamePadding}} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{heading .Title}}{{range $cmds}}{{if and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help"))}}
  {{command .Name .NamePadding}} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

{{heading "Additional Commands:"}}{{range $cmds}}{{if and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help"))}}
  {{command .Name .NamePadding}} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

{{heading "Flags:"}}
{{.LocalFlags.FlagUsages | trim}}{{end}}{{if .HasAvailableInheritedFlags}}

{{heading "Global Flags:"}}
{{.InheritedFlags.FlagUsages | trim}}{{end}}{{if .HasHelpSubCommands}}

{{heading "Additional help topics:"}}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{command .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

{{muted (printf "Use %q for more information about a command." (printf "%s [command] --help" .CommandPath))}}{{end}}
`

func configureHelp(root *cobra.Command) {
	root.SetUsageFunc(func(cmd *cobra.Command) error {
		return renderUsage(cmd.OutOrStderr(), cmd)
	})
	root.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
		w := cmd.OutOrStdout()
		description := cmd.Long
		if description == "" {
			description = cmd.Short
		}
		if description != "" {
			fmt.Fprintln(w, strings.TrimRight(description, " \t\r\n"))
			fmt.Fprintln(w)
		}
		if cmd.Runnable() || cmd.HasSubCommands() {
			// Do not use UsageString here: Cobra renders it into a buffer,
			// hiding the actual terminal destination from ui's color detection.
			if err := renderUsage(w, cmd); err != nil {
				cmd.PrintErrln(err)
			}
		}
	})
}

func renderUsage(w io.Writer, cmd *cobra.Command) error {
	// Merge persistent flags even when Help/Usage is called without Execute.
	cmd.LocalFlags()
	funcs := template.FuncMap{
		"heading": func(s string) string { return ui.Header.RenderTo(w, s) },
		"info":    func(s string) string { return ui.Info.RenderTo(w, s) },
		"muted":   func(s string) string { return ui.Muted.RenderTo(w, s) },
		"trim":    func(s string) string { return strings.TrimRight(s, " \t\r\n") },
		"command": func(name string, width int) string {
			// Pad before styling so ANSI escape sequences cannot skew columns.
			return ui.Info.RenderTo(w, fmt.Sprintf("%-*s", width, name))
		},
		"examples": func(block string) string {
			lines := strings.Split(block, "\n")
			for i, line := range lines {
				style := ui.Info
				if strings.HasPrefix(strings.TrimSpace(line), "#") {
					style = ui.Muted
				}
				lines[i] = style.RenderTo(w, line)
			}
			return strings.Join(lines, "\n")
		},
	}
	tmpl, err := template.New("usage").Funcs(funcs).Parse(usageTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, cmd)
}
