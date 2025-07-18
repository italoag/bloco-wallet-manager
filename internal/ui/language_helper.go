package ui

import (
	"strings"
)

// sanitizeNetworkKey sanitizes a key to ensure it's valid for TOML
func sanitizeNetworkKey(key string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, key)
}
