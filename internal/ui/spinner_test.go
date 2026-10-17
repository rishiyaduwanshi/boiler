package ui

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSpinnerSkipsNonTerminalOutput(t *testing.T) {
	var out bytes.Buffer
	stop := StartSpinner(&out, "Downloading...")
	stop()
	stop()
	if out.Len() != 0 {
		t.Fatalf("non-terminal spinner output = %q", out.String())
	}
}

// Only the spinner writes the buffer; the test reads it after stop has joined
// the goroutine. Notifications let us test ticks without timing-based sleeps.
type notifyingWriter struct {
	bytes.Buffer
	writes chan struct{}
}

func (w *notifyingWriter) Write(p []byte) (int, error) {
	n, err := w.Buffer.Write(p)
	select {
	case w.writes <- struct{}{}:
	default:
	}
	return n, err
}

func TestSpinnerTicksAndStops(t *testing.T) {
	out := &notifyingWriter{writes: make(chan struct{}, 16)}
	stop := startSpinner(out, "Downloading...", time.Millisecond)
	defer stop()
	for i := 0; i < 2; i++ {
		select {
		case <-out.writes:
		case <-time.After(5 * time.Second):
			t.Fatal("spinner did not render a frame")
		}
	}

	// Stopping is safe more than once, including from concurrent callers.
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); stop() }()
	}
	wg.Wait()
	got := out.String()
	if strings.Count(got, "Downloading...") < 2 {
		t.Fatalf("missing progress frames: %q", got)
	}
	if !strings.HasSuffix(got, "\r\x1b[2K") {
		t.Fatalf("spinner did not clear its line: %q", got)
	}
	if strings.Contains(got, "success") || strings.Contains(got, "✓") {
		t.Fatalf("stopping must not claim the operation succeeded: %q", got)
	}
}
