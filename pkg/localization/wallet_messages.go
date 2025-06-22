package localization

// DefaultWalletMessages returns default wallet-related messages in English
func DefaultWalletMessages() map[string]string {
	return map[string]string{
		// Error messages for wallet operations
		"error_renaming_wallet_file":        "Error renaming the wallet file: %v",
		"error_invalid_mnemonic_phrase":     "Invalid mnemonic phrase",
		"error_invalid_private_key_format":  "Invalid private key format",
		"error_invalid_private_key":         "Invalid private key: %v",
		"error_generating_mnemonic":         "Error generating mnemonic: %v",
		"error_reading_keystore_file":       "Error reading the keystore file: %v",
		"error_incorrect_keystore_password": "Incorrect password for keystore file",
		"error_getting_home_directory":      "Error getting user home directory: %v",
		"error_creating_destination_file":   "Error creating destination file: %v",
		"error_writing_destination_file":    "Error writing to destination file: %v",
		"error_reading_wallet_file":         "Error reading the wallet file: %v",
		"error_incorrect_password":          "Incorrect password",
		"error_remove_keystore_file":        "Failed to remove keystore file: %v",
	}
}

// DefaultWalletMessagesPortuguese returns wallet-related messages in Portuguese
func DefaultWalletMessagesPortuguese() map[string]string {
	return map[string]string{
		// Error messages for wallet operations in Portuguese
		"error_renaming_wallet_file":        "Erro ao renomear o arquivo da carteira: %v",
		"error_invalid_mnemonic_phrase":     "Frase mnemônica inválida",
		"error_invalid_private_key_format":  "Formato de chave privada inválido",
		"error_invalid_private_key":         "Chave privada inválida: %v",
		"error_generating_mnemonic":         "Erro ao gerar frase mnemônica: %v",
		"error_reading_keystore_file":       "Erro ao ler o arquivo keystore: %v",
		"error_incorrect_keystore_password": "Senha incorreta para o arquivo keystore",
		"error_getting_home_directory":      "Erro ao obter o diretório home do usuário: %v",
		"error_creating_destination_file":   "Erro ao criar arquivo de destino: %v",
		"error_writing_destination_file":    "Erro ao escrever no arquivo de destino: %v",
		"error_reading_wallet_file":         "Erro ao ler o arquivo da carteira: %v",
		"error_incorrect_password":          "Senha incorreta",
		"error_remove_keystore_file":        "Falha ao remover arquivo keystore: %v",
	}
}

// DefaultWalletMessagesSpanish returns wallet-related messages in Spanish
func DefaultWalletMessagesSpanish() map[string]string {
	return map[string]string{
		// Error messages for wallet operations in Spanish
		"error_renaming_wallet_file":        "Error al renombrar el archivo de la cartera: %v",
		"error_invalid_mnemonic_phrase":     "Frase mnemónica inválida",
		"error_invalid_private_key_format":  "Formato de clave privada inválido",
		"error_invalid_private_key":         "Clave privada inválida: %v",
		"error_generating_mnemonic":         "Error al generar frase mnemónica: %v",
		"error_reading_keystore_file":       "Error al leer el archivo keystore: %v",
		"error_incorrect_keystore_password": "Contraseña incorrecta para el archivo keystore",
		"error_getting_home_directory":      "Error al obtener el directorio home del usuario: %v",
		"error_creating_destination_file":   "Error al crear archivo de destino: %v",
		"error_writing_destination_file":    "Error al escribir en el archivo de destino: %v",
		"error_reading_wallet_file":         "Error al leer el archivo de la cartera: %v",
		"error_incorrect_password":          "Contraseña incorrecta",
		"error_remove_keystore_file":        "Error al eliminar archivo keystore: %v",
	}
}