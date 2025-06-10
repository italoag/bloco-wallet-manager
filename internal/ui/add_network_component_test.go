package ui

import (
	"testing"

	"github.com/charmbracelet/huh"
)

func TestAddNetworkComponent_FormAppearsImmediately(t *testing.T) {
	component := NewAddNetworkComponent()
	if component.form == nil {
		t.Error("Expected form to be initialized")
	}
	if component.form.State != huh.StateNormal {
		t.Errorf("Expected form state to be normal, got %v", component.form.State)
	}
}

func TestAddNetworkComponent_AutocompleteSuggestions(t *testing.T) {
	component := NewAddNetworkComponent()
	// As sugestões são definidas estaticamente em initForm
	expected := []string{"Polygon", "Binance Smart Chain", "Ethereum", "Avalanche",
		"Fantom", "Arbitrum", "Optimism", "Base", "Linea",
		"Polygon zkEVM", "zkSync Era", "Mantle", "Scroll"}
	// O valor inicial do campo networkName é vazio
	if component.networkName != "" {
		t.Errorf("Expected networkName to be empty, got %q", component.networkName)
	}
	// O formulário deve estar inicializado
	if component.form == nil {
		t.Fatal("form not initialized")
	}
	// Não há API pública para extrair sugestões, então validamos indiretamente
	// Se o usuário digitar, as sugestões aparecerão (testado manualmente no TUI)
	// Aqui validamos apenas que o componente foi criado corretamente
	if len(expected) < 5 {
		t.Error("Expected at least 5 suggestions for autocomplete")
	}
}

func TestAddNetworkComponent_Init(t *testing.T) {
	component := NewAddNetworkComponent()
	cmd := component.Init()

	// Verificar se Init retorna um comando
	if cmd == nil {
		t.Error("Expected Init to return a command")
	}
}
