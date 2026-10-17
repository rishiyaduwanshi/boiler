package utils

import (
	"os"
	"strings"
	"testing"
)

func TestLoggerKeepsFilesPlain(t *testing.T) {
	t.Setenv("CLICOLOR_FORCE", "1")
	logger, err := NewLogger(t.TempDir(), false)
	if err != nil {
		t.Fatal(err)
	}
	for _, level := range []LogLevel{LogLevelInfo, LogLevelWarn, LogLevelError, LogLevelDebug} {
		if err := logger.Log(level, "message"); err != nil {
			t.Fatal(err)
		}
	}
	data, err := os.ReadFile(logger.logPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "\x1b") {
		t.Fatalf("log contains ANSI sequences: %q", data)
	}
	for _, level := range []LogLevel{LogLevelInfo, LogLevelWarn, LogLevelError, LogLevelDebug} {
		if !strings.Contains(string(data), string(level)+": message\n") {
			t.Fatalf("missing %s log record: %q", level, data)
		}
	}
}
