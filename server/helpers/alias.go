package helpers

import (
	"regexp"
	"strings"
)

// GenerateAlias creates a URL-safe alias from a username
// Converts to lowercase, replaces non-alphanumeric with hyphens, removes consecutive hyphens
func GenerateAlias(username string) string {
	// Convert to lowercase
	alias := strings.ToLower(username)
	
	// Replace non-alphanumeric characters with hyphens
	re := regexp.MustCompile(`[^a-z0-9]+`)
	alias = re.ReplaceAllString(alias, "-")
	
	// Remove leading/trailing hyphens
	alias = strings.Trim(alias, "-")
	
	// Replace multiple consecutive hyphens with single hyphen
	re2 := regexp.MustCompile(`-+`)
	alias = re2.ReplaceAllString(alias, "-")
	
	// If empty after cleanup, use a default
	if alias == "" {
		alias = "user"
	}
	
	return alias
}