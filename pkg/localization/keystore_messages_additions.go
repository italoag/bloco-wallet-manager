package localization

// AddKeystoreValidationMessages adds keystore validation messages to the existing message maps
func AddKeystoreValidationMessages() {
	// Add English messages
	englishMessages := map[string]string{
		// Keystore file validation feedback
		"keystore_file_valid":     "✓ Valid keystore file detected",
		"keystore_file_not_found": "✗ File not found at the specified path",
		"keystore_access_error":   "✗ Error accessing file",
		"keystore_is_directory":   "✗ Path points to a directory, not a file",
		"keystore_not_json":       "✗ File is not a JSON file",

		// Recovery suggestions
		"keystore_recovery_file_not_found":     "Try checking the file path and ensure the file exists",
		"keystore_recovery_invalid_json":       "The file may be corrupted. Try with a different keystore file",
		"keystore_recovery_invalid_structure":  "Make sure this is a valid Ethereum keystore v3 file",
		"keystore_recovery_incorrect_password": "Please try again with the correct password",
		"keystore_recovery_general":            "Please check the file and try again",
	}

	// Add Portuguese messages
	portugueseMessages := map[string]string{
		// Keystore file validation feedback
		"keystore_file_valid":     "✓ Arquivo keystore válido detectado",
		"keystore_file_not_found": "✗ Arquivo não encontrado no caminho especificado",
		"keystore_access_error":   "✗ Erro ao acessar o arquivo",
		"keystore_is_directory":   "✗ O caminho aponta para um diretório, não um arquivo",
		"keystore_not_json":       "✗ O arquivo não é um arquivo JSON",

		// Recovery suggestions
		"keystore_recovery_file_not_found":     "Verifique o caminho do arquivo e certifique-se de que ele existe",
		"keystore_recovery_invalid_json":       "O arquivo pode estar corrompido. Tente com um arquivo keystore diferente",
		"keystore_recovery_invalid_structure":  "Certifique-se de que este é um arquivo keystore v3 Ethereum válido",
		"keystore_recovery_incorrect_password": "Por favor, tente novamente com a senha correta",
		"keystore_recovery_general":            "Por favor, verifique o arquivo e tente novamente",
	}

	// Add Spanish messages
	spanishMessages := map[string]string{
		// Keystore file validation feedback
		"keystore_file_valid":     "✓ Archivo keystore válido detectado",
		"keystore_file_not_found": "✗ Archivo no encontrado en la ruta especificada",
		"keystore_access_error":   "✗ Error al acceder al archivo",
		"keystore_is_directory":   "✗ La ruta apunta a un directorio, no a un archivo",
		"keystore_not_json":       "✗ El archivo no es un archivo JSON",

		// Recovery suggestions
		"keystore_recovery_file_not_found":     "Verifique la ruta del archivo y asegúrese de que existe",
		"keystore_recovery_invalid_json":       "El archivo puede estar dañado. Intente con un archivo keystore diferente",
		"keystore_recovery_invalid_structure":  "Asegúrese de que este es un archivo keystore v3 de Ethereum válido",
		"keystore_recovery_incorrect_password": "Por favor, intente nuevamente con la contraseña correcta",
		"keystore_recovery_general":            "Por favor, verifique el archivo e intente nuevamente",
	}

	// Add to global Labels map
	for key, value := range englishMessages {
		Labels[key] = value
	}

	// Add Portuguese and Spanish messages based on current language
	currentLang := GetCurrentLanguage()
	if currentLang == "pt" {
		for key, value := range portugueseMessages {
			Labels[key] = value
		}
	} else if currentLang == "es" {
		for key, value := range spanishMessages {
			Labels[key] = value
		}
	}
}
