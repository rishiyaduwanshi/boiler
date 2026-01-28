package utils

import (
	"fmt"

	"github.com/rishiyaduwanshi/boiler/internal/store"
)

// LoadStore creates and loads the store instance
func LoadStore(storePath string) (*store.Store, error) {
	st := store.NewStore(storePath)
	if err := st.Load(); err != nil {
		return nil, fmt.Errorf("failed to load store: %w", err)
	}
	return st, nil
}

// ConfirmAction prompts user for yes/no confirmation
func ConfirmAction(message string) bool {
	fmt.Print(message)
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}

// ParseResourceName parses resource and returns full name with version and extension
func ParseResourceName(resource string) string {
	name, version, ext := store.ParseResourceName(resource)
	fullName := name
	if version != "" {
		fullName = name + "@" + version
	}
	if ext != "" {
		fullName += ext
	}
	return fullName
}
