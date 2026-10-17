package cli

import (
	"os"
	"reflect"
	"testing"

	aliascmd "github.com/rishiyaduwanshi/boiler/internal/cli/alias"
	"github.com/rishiyaduwanshi/boiler/internal/config"
	"github.com/rishiyaduwanshi/boiler/internal/ui"
)

func TestConfigureOutputBeforeCommandDiscovery(t *testing.T) {
	old := ui.NoColor
	oldNoColor, hadNoColor := os.LookupEnv("NO_COLOR")
	t.Cleanup(func() {
		ui.NoColor = old
		if hadNoColor {
			_ = os.Setenv("NO_COLOR", oldNoColor)
		} else {
			_ = os.Unsetenv("NO_COLOR")
		}
	})
	for _, tc := range []struct {
		name string
		args []string
		want bool
	}{
		{name: "default"},
		{name: "unknown command", args: []string{"--no-color", "missing"}, want: true},
		{name: "unknown flag", args: []string{"ls", "--invalid", "--no-color"}, want: true},
		{name: "explicit true", args: []string{"ls", "--no-color=true"}, want: true},
		{name: "explicit false", args: []string{"--no-color", "--no-color=false"}},
		{name: "help", args: []string{"--help", "--no-color"}, want: true},
		{name: "short help", args: []string{"-h", "--no-color"}, want: true},
		{name: "positional argument", args: []string{"ls", "--", "--no-color"}},
		{name: "assignment value", args: []string{"var", "example=--no-color"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ConfigureOutput(tc.args)
			if ui.NoColor != tc.want {
				t.Fatalf("NoColor = %v, want %v", ui.NoColor, tc.want)
			}
		})
	}
}

func TestConfigureOutputDoesNotPropagateUntilApplied(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	old := ui.NoColor
	ui.NoColor = false
	t.Cleanup(func() { ui.NoColor = old })

	ConfigureOutput([]string{"--no-color"})

	if !ui.NoColor {
		t.Fatal("ConfigureOutput did not detect --no-color")
	}
	if got := os.Getenv("NO_COLOR"); got != "" {
		t.Fatalf("early ConfigureOutput set NO_COLOR = %q, want empty", got)
	}

	ApplyOutputToEnvironment()
	if got := os.Getenv("NO_COLOR"); got != "1" {
		t.Fatalf("NO_COLOR = %q, want %q after applying output", got, "1")
	}
}

func TestApplyOutputToEnvironmentPreservesInheritedNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	old := ui.NoColor
	ui.NoColor = false
	t.Cleanup(func() { ui.NoColor = old })

	ApplyOutputToEnvironment()
	if got := os.Getenv("NO_COLOR"); got != "1" {
		t.Fatalf("NO_COLOR = %q, want inherited value %q", got, "1")
	}
}

func TestApplyOutputToEnvironmentPreservesInheritedValueAfterAliasExpansion(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	old := ui.NoColor
	ui.NoColor = false
	t.Cleanup(func() { ui.NoColor = old })

	cfg := config.DefaultConfig()
	cfg.Aliases["find"] = "search --"
	aliascmd.Setup(cfg, nil)
	t.Cleanup(func() { aliascmd.Setup(nil, nil) })

	original := []string{"find", "--no-color"}
	ConfigureOutput(original)
	if !ui.NoColor {
		t.Fatal("early scan did not detect --no-color")
	}

	expanded := aliascmd.ExpandFirstCommandAlias(original)
	wantExpanded := []string{"search", "--", "--no-color"}
	if !reflect.DeepEqual(expanded, wantExpanded) {
		t.Fatalf("ExpandFirstCommandAlias() = %v, want %v", expanded, wantExpanded)
	}

	ConfigureOutput(expanded)
	if ui.NoColor {
		t.Fatal("final scan treated --no-color after -- as a flag")
	}
	ApplyOutputToEnvironment()
	if got := os.Getenv("NO_COLOR"); got != "1" {
		t.Fatalf("NO_COLOR = %q, want inherited value %q after alias expansion", got, "1")
	}
}
