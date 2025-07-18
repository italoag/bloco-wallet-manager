package localization

// currentLanguage stores the current language code
var currentLanguage string = "en" // Default to English

// GetCurrentLanguage returns the current language code
func GetCurrentLanguage() string {
	return currentLanguage
}

// SetCurrentLanguage sets the current language code
func SetCurrentLanguage(lang string) {
	currentLanguage = lang
}
