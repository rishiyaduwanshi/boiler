package remote

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadsWithProgress(t *testing.T) {
	for _, tc := range []struct {
		name      string
		status    int
		chunked   bool
		truncated bool
		wantErr   string
	}{
		{name: "known length", status: http.StatusOK},
		{name: "unknown length", status: http.StatusOK, chunked: true},
		{name: "not found", status: http.StatusNotFound, wantErr: "HTTP 404"},
		{name: "server error", status: http.StatusInternalServerError, wantErr: "HTTP 500"},
		{name: "truncated body", status: http.StatusOK, truncated: true, wantErr: "unexpected EOF"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			const body = "snippet contents\n"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.truncated {
					w.Header().Set("Content-Length", "1000")
				}
				w.WriteHeader(tc.status)
				if tc.chunked {
					w.(http.Flusher).Flush()
				}
				_, _ = fmt.Fprint(w, body)
			}))
			defer server.Close()

			data, err := downloadFile(server.URL)
			assertDownloadError(t, err, tc.wantErr)
			if tc.wantErr == "" && string(data) != body {
				t.Fatalf("downloaded data = %q, want %q", data, body)
			}

			path := filepath.Join(t.TempDir(), "stack.zip")
			err = downloadToFile(server.URL, path)
			assertDownloadError(t, err, tc.wantErr)
			if tc.wantErr == "" {
				data, err = os.ReadFile(path)
				if err != nil || string(data) != body {
					t.Fatalf("downloaded file = %q, error = %v", data, err)
				}
			}
		})
	}
}

func TestDownloadsRejectInvalidURL(t *testing.T) {
	_, err := downloadFile(":invalid")
	assertDownloadError(t, err, "failed to create request")
	err = downloadToFile(":invalid", filepath.Join(t.TempDir(), "download"))
	assertDownloadError(t, err, "failed to create request")
}

func assertDownloadError(t *testing.T, err error, want string) {
	t.Helper()
	if want == "" {
		if err != nil {
			t.Fatalf("unexpected download error: %v", err)
		}
	} else if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("download error = %v, want containing %q", err, want)
	}
}
