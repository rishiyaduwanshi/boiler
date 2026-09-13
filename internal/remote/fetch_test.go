package remote

import "testing"

func TestParseHostedStackRefAndSubPath(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantRef     string
		wantSubPath string
		wantOK      bool
	}{
		{
			name:        "github tree subpath",
			input:       "https://github.com/spf13/cobra/tree/main/doc",
			wantRef:     "main",
			wantSubPath: "doc",
			wantOK:      true,
		},
		{
			name:        "github blob subpath",
			input:       "https://github.com/spf13/cobra/blob/main/doc",
			wantRef:     "main",
			wantSubPath: "doc",
			wantOK:      true,
		},
		{
			name:        "github tree root",
			input:       "https://github.com/spf13/cobra/tree/main",
			wantRef:     "main",
			wantSubPath: ".",
			wantOK:      true,
		},
		{
			name:        "gitlab tree subpath",
			input:       "https://gitlab.com/group/repo/-/tree/main/templates/api",
			wantRef:     "main",
			wantSubPath: "templates/api",
			wantOK:      true,
		},
		{
			name:        "bitbucket src subpath",
			input:       "https://bitbucket.org/team/repo/src/main/pkg",
			wantRef:     "main",
			wantSubPath: "pkg",
			wantOK:      true,
		},
		{
			name:   "regular repo url",
			input:  "https://github.com/spf13/cobra",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRef, gotSubPath, gotOK := parseHostedStackRefAndSubPath(tt.input)
			if gotOK != tt.wantOK {
				t.Fatalf("ok = %v, want %v", gotOK, tt.wantOK)
			}
			if gotRef != tt.wantRef {
				t.Fatalf("ref = %q, want %q", gotRef, tt.wantRef)
			}
			if gotSubPath != tt.wantSubPath {
				t.Fatalf("subPath = %q, want %q", gotSubPath, tt.wantSubPath)
			}
		})
	}
}

func TestAuthHeaderForURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		env       map[string]string
		wantKey   string
		wantValue string
	}{
		// ── GitHub exact hosts ────────────────────────────────────────────
		{
			name:      "github.com with token",
			url:       "https://github.com/owner/repo",
			env:       map[string]string{"GITHUB_TOKEN": "gh-tok"},
			wantKey:   "Authorization",
			wantValue: "Bearer gh-tok",
		},
		{
			name:      "api.github.com with token",
			url:       "https://api.github.com/repos/owner/repo/tarball/main",
			env:       map[string]string{"GITHUB_TOKEN": "gh-tok"},
			wantKey:   "Authorization",
			wantValue: "Bearer gh-tok",
		},
		{
			name:      "raw.githubusercontent.com with token",
			url:       "https://raw.githubusercontent.com/owner/repo/main/file.go",
			env:       map[string]string{"GITHUB_TOKEN": "gh-tok"},
			wantKey:   "Authorization",
			wantValue: "Bearer gh-tok",
		},
		{
			name: "github.com without token set",
			url:  "https://github.com/owner/repo",
			env:  map[string]string{},
		},

		// ── GitLab exact hosts ────────────────────────────────────────────
		{
			name:      "gitlab.com with token",
			url:       "https://gitlab.com/owner/repo",
			env:       map[string]string{"GITLAB_TOKEN": "gl-tok"},
			wantKey:   "Authorization",
			wantValue: "Bearer gl-tok",
		},
		{
			name: "gitlab.com without token set",
			url:  "https://gitlab.com/owner/repo",
			env:  map[string]string{},
		},

		// ── Bitbucket exact hosts ─────────────────────────────────────────
		{
			name:      "bitbucket.org with token",
			url:       "https://bitbucket.org/owner/repo",
			env:       map[string]string{"BITBUCKET_TOKEN": "bb-tok"},
			wantKey:   "Authorization",
			wantValue: "Bearer bb-tok",
		},
		{
			name: "bitbucket.org without token set",
			url:  "https://bitbucket.org/owner/repo",
			env:  map[string]string{},
		},

		// ── Deceptive / spoofed URLs → must never send a token ───────────
		{
			name: "github.com in userinfo (SSRF spoof)",
			url:  "http://github.com@attacker.example/archive.zip",
			env:  map[string]string{"GITHUB_TOKEN": "gh-tok"},
		},
		{
			name: "github.com in URL path",
			url:  "https://attacker.example/github.com/owner/repo",
			env:  map[string]string{"GITHUB_TOKEN": "gh-tok"},
		},
		{
			name: "github.com in query string",
			url:  "https://attacker.example/download?from=github.com",
			env:  map[string]string{"GITHUB_TOKEN": "gh-tok"},
		},
		{
			name: "unknown host with all tokens set",
			url:  "https://myregistry.example/stack.zip",
			env: map[string]string{
				"GITHUB_TOKEN":    "gh-tok",
				"GITLAB_TOKEN":    "gl-tok",
				"BITBUCKET_TOKEN": "bb-tok",
			},
		},
		{
			name: "invalid URL",
			url:  "://bad-url",
			env:  map[string]string{"GITHUB_TOKEN": "gh-tok"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			gotKey, gotValue := authHeaderForURL(tt.url)

			if gotKey != tt.wantKey {
				t.Errorf("key = %q, want %q", gotKey, tt.wantKey)
			}
			if gotValue != tt.wantValue {
				t.Errorf("value = %q, want %q", gotValue, tt.wantValue)
			}
		})
	}
}
