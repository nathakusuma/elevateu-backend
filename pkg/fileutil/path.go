package fileutil

import "regexp"

// IsValidPath checks if the path contains only allowed characters
func IsValidPath(path string) bool {
	// Define allowed characters pattern
	// This allows:
	// - alphanumeric characters
	// - forward slashes for directory separation
	// - common special characters like dash, underscore
	// - dots for file extensions
	pattern := regexp.MustCompile(`^[a-zA-Z0-9\-_./]+$`)
	return pattern.MatchString(path)
}
