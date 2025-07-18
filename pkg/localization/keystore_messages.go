package localization

// GetKeystoreErrorMessage returns a localized error message for a keystore error key
func GetKeystoreErrorMessage(key string) string {
	// Use Labels map directly for simplicity
	if value, ok := Labels[key]; ok {
		return value
	}
	return key
}

// FormatKeystoreErrorWithField formats a keystore error message with a field name
func FormatKeystoreErrorWithField(key string, field string) string {
	message := GetKeystoreErrorMessage(key)
	if field != "" {
		return message + " (" + field + ")"
	}
	return message
}
