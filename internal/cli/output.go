package cli

import (
	"io"
	"os"

	"github.com/rishiyaduwanshi/boiler/internal/ui"
	"github.com/spf13/pflag"
)

// ConfigureOutput applies the output preference before config loading and
// command discovery, which can fail before Cobra parses persistent flags.
// Cobra remains responsible for reporting invalid flags and argument values.
// Environment propagation is deliberately separate because aliases are
// expanded only after config loading.
func ConfigureOutput(args []string) {
	flags := pflag.NewFlagSet("output", pflag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.ParseErrorsAllowlist.UnknownFlags = true
	flags.BoolVar(&ui.NoColor, "no-color", false, "")
	// Do not let help stop this preliminary scan before a trailing --no-color.
	flags.BoolP("help", "h", false, "")
	_ = flags.Parse(args)
}

// ApplyOutputToEnvironment propagates the final output preference to child
// processes after command aliases have been expanded. Inherited NO_COLOR values
// are preserved when no explicit opt-out is requested.
func ApplyOutputToEnvironment() {
	if ui.NoColor {
		_ = os.Setenv("NO_COLOR", "1")
	}
}
