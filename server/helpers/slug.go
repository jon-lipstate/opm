package helpers

import (
	"regexp"
	"strings"
)

// GenerateSlug creates a URL-safe slug from a username
// Converts to lowercase, replaces non-alphanumeric with hyphens, removes consecutive hyphens
func GenerateUserSlug(username string) string {
	// Convert to lowercase
	slug := strings.ToLower(username)

	// Replace non-alphanumeric characters with hyphens
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	// Replace multiple consecutive hyphens with single hyphen
	re2 := regexp.MustCompile(`-+`)
	slug = re2.ReplaceAllString(slug, "-")

	// If empty after cleanup, use a default
	if slug == "" {
		slug = "user"
	}

	return slug
}
