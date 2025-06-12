package localization

// DefaultCryptoMessages retorna mensagens padrão para criptografia em inglês
func DefaultCryptoMessages() map[string]string {
	return map[string]string{
		// Mensagens para criptografia
		"error_empty_password":                 "Password cannot be empty",
		"error_generate_salt":                  "Failed to generate salt",
		"error_empty_encrypted_mnemonic":       "Encrypted mnemonic cannot be empty",
		"error_decode_mnemonic":                "Failed to decode encrypted mnemonic",
		"error_invalid_mnemonic_format":        "Invalid encrypted mnemonic format: too short",
		"error_invalid_password":               "Invalid password: hash verification failed",
		"error_crypto_service_not_initialized": "Crypto service not initialized",
	}
}

// DefaultCryptoMessagesPortuguese retorna mensagens padrão para criptografia em português
func DefaultCryptoMessagesPortuguese() map[string]string {
	return map[string]string{
		// Mensagens para criptografia em português
		"error_empty_password":                 "A senha não pode estar vazia",
		"error_generate_salt":                  "Falha ao gerar salt criptográfico",
		"error_empty_encrypted_mnemonic":       "A frase mnemônica criptografada não pode estar vazia",
		"error_decode_mnemonic":                "Falha ao decodificar a frase mnemônica criptografada",
		"error_invalid_mnemonic_format":        "Formato de frase mnemônica criptografada inválido: muito curta",
		"error_invalid_password":               "Senha inválida: verificação de hash falhou",
		"error_crypto_service_not_initialized": "Serviço de criptografia não inicializado",
	}
}

// DefaultCryptoMessagesSpanish retorna mensagens padrão para criptografia em espanhol
func DefaultCryptoMessagesSpanish() map[string]string {
	return map[string]string{
		// Mensagens para criptografia em espanhol
		"error_empty_password":                 "La contraseña no puede estar vacía",
		"error_generate_salt":                  "Error al generar salt criptográfico",
		"error_empty_encrypted_mnemonic":       "La frase mnemónica cifrada no puede estar vacía",
		"error_decode_mnemonic":                "Error al decodificar la frase mnemónica cifrada",
		"error_invalid_mnemonic_format":        "Formato de frase mnemónica cifrada inválido: demasiado corta",
		"error_invalid_password":               "Contraseña inválida: la verificación del hash falló",
		"error_crypto_service_not_initialized": "Servicio de criptografía no inicializado",
	}
}
