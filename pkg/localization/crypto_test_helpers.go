package localization

// InitCryptoMessagesForTesting inicializa as mensagens de criptografia para testes unitários
// sem depender de arquivos de idioma ou do Viper
func InitCryptoMessagesForTesting() {
	// Inicializa o mapa global se ainda não estiver inicializado
	if Labels == nil {
		Labels = make(map[string]string)
	}

	// Adicionar todas as mensagens de criptografia em inglês para uso nos testes
	for key, value := range DefaultCryptoMessages() {
		Labels[key] = value
	}
}

// GetForTesting para testes retorna a mensagem diretamente do mapa Labels
// Esta função é uma versão simplificada da função Get para uso em testes
func GetForTesting(key string) string {
	if Labels == nil {
		return key
	}

	if value, ok := Labels[key]; ok {
		return value
	}
	return key
}
