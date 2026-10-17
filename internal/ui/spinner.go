package ui

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// StartSpinner provides feedback while a synchronous operation is running.
// Defer the returned stop function immediately: it waits for the animation to
// stop and clears the line before the caller prints a success or error message.
// Pipes, dumb terminals and no-color output never receive animation sequences.
func StartSpinner(w io.Writer, message string) func() {
	if !colorEnabled(w) || lipgloss.NewRenderer(w).ColorProfile() == termenv.Ascii {
		return func() {}
	}
	return startSpinner(w, message, spinner.Dot.FPS)
}

func startSpinner(w io.Writer, message string, interval time.Duration) func() {
	model := spinner.New(spinner.WithSpinner(spinner.Dot))
	// Drive the component with a ticker rather than starting a full-screen TUI.
	// This leaves stdin, signal handling and the terminal's input mode untouched.
	render := func() {
		_, _ = fmt.Fprint(w, "\r\x1b[2K", Info.render(w, model.View()+" "+message))
	}
	render()

	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				_, _ = fmt.Fprint(w, "\r\x1b[2K")
				return
			case now := <-ticker.C:
				// The local ticker owns scheduling, so no Bubble Tea command is run.
				model, _ = model.Update(spinner.TickMsg{Time: now, ID: model.ID()})
				render()
			}
		}
	}()

	var once sync.Once
	return func() {
		once.Do(func() { close(stop) })
		<-done
	}
}
