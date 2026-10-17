// Package ui provides the shared presentation styles for Boiler's terminal output.
package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/muesli/termenv"
)

// NoColor is bound to the persistent --no-color flag. Cobra sets it before
// rendering help or running commands. Do not change it while output is active.
var NoColor bool

// Style renders semantic CLI output without changing plain-text output.
// Keeping the underlying style private ensures every render uses the output
// destination's capabilities rather than Lip Gloss's global renderer.
type Style struct {
	style lipgloss.Style
}

var (
	Error   = Style{lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)}
	Success = Style{lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)}
	Info    = Style{lipgloss.NewStyle().Foreground(lipgloss.Color("86"))}
	Warning = Style{lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)}
	Path    = Style{lipgloss.NewStyle().Foreground(lipgloss.Color("14"))}
	Header  = Style{lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("62")).Bold(true).Padding(0, 1)}
	Muted   = Style{lipgloss.NewStyle().Foreground(lipgloss.Color("245"))}
)

// colorEnabled checks the actual destination: redirecting stdout must not
// disable stderr's styling (or accidentally color a file when stderr is a TTY).
func colorEnabled(w io.Writer) bool {
	if NoColor || os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}
	f, ok := w.(interface{ Fd() uintptr })
	return ok && term.IsTerminal(f.Fd())
}

func (s Style) render(w io.Writer, text string) string {
	if !colorEnabled(w) {
		return text
	}
	renderer := lipgloss.NewRenderer(w)
	if renderer.ColorProfile() == termenv.Ascii {
		return text
	}
	body := strings.Trim(text, "\n")
	if body == "" {
		return text
	}
	start := len(text) - len(strings.TrimLeft(text, "\n"))
	end := len(strings.TrimRight(text, "\n"))
	return text[:start] + s.style.Renderer(renderer).Render(body) + text[end:]
}

// RenderTo styles text for a specific output destination without writing it.
// Use it when composing output that may go to a terminal, file, or buffer.
func (s Style) RenderTo(w io.Writer, text string) string {
	return s.render(w, text)
}

// Render styles a value destined for stdout, e.g. a path inside a message.
func (s Style) Render(text string) string {
	return s.render(os.Stdout, text)
}

func (s Style) Printf(format string, args ...any) {
	s.Fprintf(os.Stdout, format, args...)
}

func (s Style) Println(args ...any) {
	_, _ = fmt.Fprint(os.Stdout, s.render(os.Stdout, fmt.Sprintln(args...)))
}

func (s Style) Fprintf(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprint(w, s.render(w, fmt.Sprintf(format, args...)))
}
