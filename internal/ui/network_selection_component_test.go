package ui

import (
	"blocowallet/pkg/config"
	"testing"

	"github.com/charmbracelet/huh"
)

func TestNetworkSelectionComponent_FormAppearsImmediately(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = &config.Config{} // Usar configuração vazia para testes se falhar o carregamento
	}
	component := NewNetworkSelectionComponent(cfg)
	if component.form == nil {
		t.Error("Expected form to be initialized")
	}
	if component.form.State != huh.StateNormal {
		t.Errorf("Expected form state to be normal, got %v", component.form.State)
	}
}

func TestNetworkSelectionComponent_Init(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = &config.Config{} // Usar configuração vazia para testes se falhar o carregamento
	}
	component := NewNetworkSelectionComponent(cfg)
	cmd := component.Init()

	// Verificar se Init retorna um comando
	if cmd == nil {
		t.Error("Expected Init to return a command")
	}
}
