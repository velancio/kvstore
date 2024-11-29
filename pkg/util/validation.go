package util

import (
	"fmt"
	"regexp"
)

// ValidateKvPair validates a key-value pair.
func ValidateKvPair(key string, value string) error {
	// Key validation
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	// Value validation
	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	// Additional custom validations
	if containsInvalidChars(key) {
		return fmt.Errorf("key contains invalid characters")
	}

	return nil
}

// containsInvalidChars checks if the given string contains invalid characters.
func containsInvalidChars(s string) bool {
	// Define a regular expression to match invalid characters
	invalidChars := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return invalidChars.MatchString(s)
}
