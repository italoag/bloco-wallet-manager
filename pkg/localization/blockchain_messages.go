package localization

// DefaultBlockchainMessages returns default blockchain-related messages in English
func DefaultBlockchainMessages() map[string]string {
	return map[string]string{
		// Error messages for blockchain operations
		"error_failed_get_balance": "Failed to get balance on %s: %v",
	}
}

// DefaultBlockchainMessagesPortuguese returns blockchain-related messages in Portuguese
func DefaultBlockchainMessagesPortuguese() map[string]string {
	return map[string]string{
		// Error messages for blockchain operations in Portuguese
		"error_failed_get_balance": "Falha ao obter saldo na rede %s: %v",
	}
}

// DefaultBlockchainMessagesSpanish returns blockchain-related messages in Spanish
func DefaultBlockchainMessagesSpanish() map[string]string {
	return map[string]string{
		// Error messages for blockchain operations in Spanish
		"error_failed_get_balance": "Error al obtener saldo en la red %s: %v",
	}
}