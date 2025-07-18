package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ConfigFileManager gerencia operações de leitura e escrita no arquivo de configuração
type ConfigFileManager struct {
	configPath string
}

// NewConfigFileManager cria uma nova instância do gerenciador
func NewConfigFileManager(configPath string) *ConfigFileManager {
	return &ConfigFileManager{
		configPath: configPath,
	}
}

// ReadConfig lê o arquivo de configuração atual
func (cfm *ConfigFileManager) ReadConfig() ([]string, error) {
	// Verificar se o arquivo existe
	if _, err := os.Stat(cfm.configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("arquivo de configuração não encontrado: %w", err)
	}

	// Ler o conteúdo do arquivo
	content, err := os.ReadFile(cfm.configPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler arquivo de configuração: %w", err)
	}

	// Dividir o conteúdo em linhas
	lines := strings.Split(string(content), "\n")
	return lines, nil
}

// WriteConfig escreve o conteúdo atualizado no arquivo de configuração
func (cfm *ConfigFileManager) WriteConfig(lines []string) error {
	// Criar um backup antes de escrever
	backupPath, err := cfm.BackupConfig()
	if err != nil {
		return fmt.Errorf("falha ao criar backup: %w", err)
	}

	// Juntar as linhas em uma string
	content := strings.Join(lines, "\n")

	// Escrever para um arquivo temporário primeiro
	tempPath := cfm.configPath + ".tmp"
	err = os.WriteFile(tempPath, []byte(content), 0644)
	if err != nil {
		// Tentar restaurar o backup em caso de erro
		_ = cfm.RestoreConfig(backupPath)
		return fmt.Errorf("falha ao escrever arquivo temporário: %w", err)
	}

	// Renomear o arquivo temporário para o arquivo de configuração
	err = os.Rename(tempPath, cfm.configPath)
	if err != nil {
		// Tentar restaurar o backup em caso de erro
		_ = cfm.RestoreConfig(backupPath)
		return fmt.Errorf("falha ao renomear arquivo temporário: %w", err)
	}

	return nil
}

// BackupConfig cria um backup do arquivo de configuração
func (cfm *ConfigFileManager) BackupConfig() (string, error) {
	// Verificar se o arquivo existe
	if _, err := os.Stat(cfm.configPath); os.IsNotExist(err) {
		return "", fmt.Errorf("arquivo de configuração não encontrado: %w", err)
	}

	// Criar nome do arquivo de backup com timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s.bak", cfm.configPath, timestamp)

	// Copiar o arquivo para o backup
	content, err := os.ReadFile(cfm.configPath)
	if err != nil {
		return "", fmt.Errorf("falha ao ler arquivo para backup: %w", err)
	}

	err = os.WriteFile(backupPath, content, 0644)
	if err != nil {
		return "", fmt.Errorf("falha ao escrever arquivo de backup: %w", err)
	}

	return backupPath, nil
}

// RestoreConfig restaura o arquivo de configuração a partir de um backup
func (cfm *ConfigFileManager) RestoreConfig(backupPath string) error {
	// Verificar se o backup existe
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("arquivo de backup não encontrado: %w", err)
	}

	// Ler o conteúdo do backup
	content, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("falha ao ler arquivo de backup: %w", err)
	}

	// Escrever o conteúdo do backup no arquivo de configuração
	err = os.WriteFile(cfm.configPath, content, 0644)
	if err != nil {
		return fmt.Errorf("falha ao restaurar backup: %w", err)
	}

	return nil
}

// EnsureConfigDir garante que o diretório de configuração existe
func (cfm *ConfigFileManager) EnsureConfigDir() error {
	dir := filepath.Dir(cfm.configPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("falha ao criar diretório de configuração: %w", err)
		}
	}
	return nil
}

// ValidateConfigFile verifica se o arquivo de configuração é válido
func (cfm *ConfigFileManager) ValidateConfigFile() error {
	// Aqui poderíamos implementar uma validação mais completa do arquivo TOML
	// Por enquanto, apenas verificamos se o arquivo existe e pode ser lido
	_, err := cfm.ReadConfig()
	return err
}
