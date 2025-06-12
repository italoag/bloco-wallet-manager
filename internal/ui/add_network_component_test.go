package ui

import (
	"testing"
)

func TestAddNetworkComponent_FormAppearsImmediately(t *testing.T) {
	component := NewAddNetworkComponent()
	if component.searchInput.Value() != "" {
		t.Error("Expected search input to be empty initially")
	}
	if component.nameInput.Value() != "" {
		t.Error("Expected name input to be empty initially")
	}
	if component.chainIDInput.Value() != "" {
		t.Error("Expected chain ID input to be empty initially")
	}
	if component.symbolInput.Value() != "" {
		t.Error("Expected symbol input to be empty initially")
	}
	if component.rpcEndpointInput.Value() != "" {
		t.Error("Expected RPC endpoint input to be empty initially")
	}
}

func TestAddNetworkComponent_AutocompleteSuggestions(t *testing.T) {
	component := NewAddNetworkComponent()
	// As sugestões são carregadas dinamicamente através do chainListService
	// Vamos testar se os inputs estão configurados corretamente para sugestões

	// O valor inicial do campo search deve ser vazio
	if component.searchInput.Value() != "" {
		t.Errorf("Expected search input to be empty, got %q", component.searchInput.Value())
	}

	// O input de pesquisa deve estar configurado para mostrar sugestões
	if !component.searchInput.ShowSuggestions {
		t.Error("Expected searchInput to show suggestions")
	}

	// Verificar se o chainListService foi inicializado
	if component.chainListService == nil {
		t.Error("Expected chainListService to be initialized")
	}
}

func TestAddNetworkComponent_Init(t *testing.T) {
	component := NewAddNetworkComponent()
	cmd := component.Init()

	// Verificar se Init retorna um comando
	if cmd == nil {
		t.Error("Expected Init to return a command")
	}

	// Verificar se o foco inicial está no campo de pesquisa
	if !component.searchInput.Focused() {
		t.Error("Expected searchInput to be focused initially")
	}

	// Verificar se isSearchFocused está definido corretamente
	if !component.isSearchFocused {
		t.Error("Expected isSearchFocused to be true initially")
	}
}
