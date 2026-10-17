package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rishiyaduwanshi/boiler/internal/ui"
)

func Prompt(message string) (string, error) {
	ui.Info.Printf("%s", message)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func PromptWithDefault(message, defaultValue string) (string, error) {
	promptMsg := fmt.Sprintf("%s [%s]: ", message, defaultValue)
	input, err := Prompt(promptMsg)
	if err != nil {
		return "", err
	}
	if input == "" {
		return defaultValue, nil
	}
	return input, nil
}

// PromptString is a convenience wrapper for PromptWithDefault
func PromptString(message, defaultValue string) string {
	result, err := PromptWithDefault(message, defaultValue)
	if err != nil {
		return defaultValue
	}
	return result
}

func PromptYesNo(message string) (bool, error) {
	input, err := Prompt(message + " (y/n): ")
	if err != nil {
		return false, err
	}
	input = strings.ToLower(input)
	return input == "y" || input == "yes", nil
}
