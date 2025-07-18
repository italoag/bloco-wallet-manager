# Keystore Validation Localization

This document describes the localization system for keystore validation error messages in the BLOCO Wallet Manager.

## Overview

The keystore validation system provides detailed error messages in multiple languages to help users understand and resolve issues when importing keystore files. The system supports the following languages:

- English (default)
- Portuguese (pt)
- Spanish (es)

## Message Categories

The localization system includes the following categories of messages:

### Validation Feedback

These messages provide immediate feedback during the keystore validation process:

- `keystore_file_valid`: Confirmation that a valid keystore file was detected
- `keystore_file_not_found`: Error when the specified file path doesn't exist
- `keystore_access_error`: Error when the file exists but cannot be accessed
- `keystore_is_directory`: Error when the path points to a directory instead of a file
- `keystore_not_json`: Error when the file is not a valid JSON file

### Recovery Suggestions

These messages provide actionable suggestions to help users recover from errors:

- `keystore_recovery_file_not_found`: Suggestion when a file is not found
- `keystore_recovery_invalid_json`: Suggestion when the JSON is invalid
- `keystore_recovery_invalid_structure`: Suggestion when the keystore structure is invalid
- `keystore_recovery_incorrect_password`: Suggestion when the password is incorrect
- `keystore_recovery_general`: General recovery suggestion

## Implementation

The localization system is implemented in the following files:

- `pkg/localization/keystore_messages.go`: Core functions for retrieving localized messages
- `pkg/localization/keystore_messages_additions.go`: Language-specific message definitions
- `pkg/localization/current_language.go`: Language selection management

### Language Selection

The current language is stored in the `currentLanguage` variable in `pkg/localization/current_language.go`. The default language is English ("en").

To change the language:

```go
localization.SetCurrentLanguage("pt") // Set to Portuguese
localization.SetCurrentLanguage("es") // Set to Spanish
localization.SetCurrentLanguage("en") // Set to English
```

### Message Initialization

Messages are initialized when the application starts by calling:

```go
localization.AddKeystoreValidationMessages()
```

This function adds all keystore validation messages to the global `Labels` map based on the current language setting.

## Usage in Code

To use localized keystore error messages in code:

```go
// Get a simple error message
errorMsg := localization.GetKeystoreErrorMessage("keystore_file_not_found")

// Format an error message with a field name
formattedMsg := localization.FormatKeystoreErrorWithField("keystore_not_json", "file.json")
```

## Testing

To test the localization system with different languages:

```go
// Set language to Portuguese
localization.SetCurrentLanguage("pt")
localization.AddKeystoreValidationMessages()

// Test a message
msg := localization.GetKeystoreErrorMessage("keystore_file_not_found")
// msg should be "✗ Arquivo não encontrado no caminho especificado"

// Set language to Spanish
localization.SetCurrentLanguage("es")
localization.AddKeystoreValidationMessages()

// Test a message
msg = localization.GetKeystoreErrorMessage("keystore_file_not_found")
// msg should be "✗ Archivo no encontrado en la ruta especificada"
```

## Adding New Languages

To add support for a new language:

1. Add a new map in `pkg/localization/keystore_messages_additions.go`
2. Add the language code to the conditional logic in `AddKeystoreValidationMessages()`
3. Ensure all message keys have translations in the new language