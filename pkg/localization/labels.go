package localization

import (
	"blocowallet/pkg/config"
)

// Labels is a map of localized strings
// This is kept for backward compatibility
// The actual implementation has been moved to locale.go
var Labels map[string]string

// SetLanguage initializes the localization system
// This function is kept for backward compatibility
// It delegates to InitLocalization in locale.go
func SetLanguage(lang string, appDir string) error {
	// Set the current language
	SetCurrentLanguage(lang)

	// Create a minimal config for InitLocalization
	cfg := &config.Config{
		Language:  lang,
		LocaleDir: appDir + "/locale",
	}

	// Delegate to the new implementation
	return InitLocalization(cfg)
}
