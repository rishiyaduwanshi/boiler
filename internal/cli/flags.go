package cli

// Flag constants for consistent documentation across commands
const (
	// Short flags
	FlagSnippetsShort = "n" // -n for snippets
	FlagStacksShort   = "k" // -k for stacks
	FlagForceShort    = "f" // -f for force operations
	FlagAllShort      = "a" // -a for all resources

	// Long flags
	FlagSnippets = "snippets"
	FlagStacks   = "stacks"
	FlagForce    = "force"
	FlagAll      = "all"

	// Flag descriptions
	DescSnippetsOnly = "Snippets only"
	DescStacksOnly   = "Stacks only"
	DescForce        = "Force operation without confirmation"
	DescCleanAll     = "Clean all resources"
)
