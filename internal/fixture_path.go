package internal

import (
	"path/filepath"
	"regexp"
)

// FixturePath returns the path to a fixture for a particular event type.
// Non alphanumeric characters are replaced with dashes to avoid
// having ':' characters in filenames (which might cause weirdness in Windows).
//
// Providing 'cf:request.created' as the event returns 'cf-request-created'.
func FixturePath(event string) string {
	// replace cf:request.created with cf-request-created
	sanitized := regexp.MustCompile(`[^a-zA-Z\d]`).ReplaceAllString(event, "-")
	f := sanitized + ".json"
	return filepath.Join("fixtures", f)
}
