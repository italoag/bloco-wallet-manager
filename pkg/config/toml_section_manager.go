package config

import (
	"fmt"
	"strings"
)

// TOMLSectionManager gerencia seções em arquivos TOML
type TOMLSectionManager struct{}

// NewTOMLSectionManager cria uma nova instância do gerenciador de seções
func NewTOMLSectionManager() *TOMLSectionManager {
	return &TOMLSectionManager{}
}

// FindSection encontra o início e fim de uma seção específica
// Retorna o índice de início da seção e o índice de fim (exclusivo)
// Se a seção não for encontrada, retorna -1, -1
func (tsm *TOMLSectionManager) FindSection(lines []string, sectionName string) (int, int) {
	sectionHeader := "[" + sectionName + "]"
	start := -1
	end := -1

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Encontrou o início da seção
		if trimmedLine == sectionHeader {
			start = i
			continue
		}

		// Se já encontramos o início da seção e encontramos outra seção, então esta é o fim
		if start != -1 && strings.HasPrefix(trimmedLine, "[") && !strings.HasPrefix(trimmedLine, "["+sectionName+".") {
			end = i
			break
		}
	}

	// Se encontramos o início mas não o fim, o fim é o final do arquivo
	if start != -1 && end == -1 {
		end = len(lines)
	}

	return start, end
}

// RemoveSection remove uma seção específica do conteúdo
func (tsm *TOMLSectionManager) RemoveSection(lines []string, sectionName string) []string {
	start, end := tsm.FindSection(lines, sectionName)

	// Se a seção não foi encontrada, retorna o conteúdo original
	if start == -1 {
		return lines
	}

	// Remove a seção
	result := append([]string{}, lines[:start]...)
	if end < len(lines) {
		result = append(result, lines[end:]...)
	}

	return result
}

// RemoveSubSections remove todas as subseções de uma seção específica
func (tsm *TOMLSectionManager) RemoveSubSections(lines []string, sectionName string) []string {
	var result []string
	inSubSection := false
	subSectionPrefix := "[" + sectionName + "."

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Verifica se estamos entrando em uma subseção
		if strings.HasPrefix(trimmedLine, subSectionPrefix) {
			inSubSection = true
			continue
		}

		// Verifica se estamos saindo de uma subseção
		if inSubSection && strings.HasPrefix(trimmedLine, "[") && !strings.HasPrefix(trimmedLine, subSectionPrefix) {
			inSubSection = false
		}

		// Adiciona a linha se não estiver em uma subseção
		if !inSubSection {
			result = append(result, line)
		}
	}

	return result
}

// AddSection adiciona uma nova seção ao conteúdo
func (tsm *TOMLSectionManager) AddSection(lines []string, sectionName string, sectionContent []string) []string {
	// Primeiro, remove qualquer seção existente com o mesmo nome
	lines = tsm.RemoveSection(lines, sectionName)

	// Adiciona uma linha em branco antes da nova seção se não terminar com uma
	if len(lines) > 0 && lines[len(lines)-1] != "" {
		lines = append(lines, "")
	}

	// Adiciona o cabeçalho da seção
	lines = append(lines, "["+sectionName+"]")

	// Adiciona o conteúdo da seção
	if len(sectionContent) > 0 {
		lines = append(lines, sectionContent...)
	}

	return lines
}

// FormatNetworkSection formata a seção de redes corretamente
func (tsm *TOMLSectionManager) FormatNetworkSection(networks map[string]Network) []string {
	var result []string

	// Se não há redes, retorna uma lista vazia
	if len(networks) == 0 {
		return result
	}

	// Adiciona a seção [networks]
	result = append(result, "[networks]")

	// Para cada rede, adiciona uma subseção
	for key, network := range networks {
		// Adiciona uma linha em branco para separar as subseções
		result = append(result, "")

		// Adiciona o cabeçalho da subseção
		result = append(result, fmt.Sprintf("[networks.%s]", key))

		// Adiciona os campos da rede
		result = append(result, fmt.Sprintf("name = %q", network.Name))
		result = append(result, fmt.Sprintf("rpc_endpoint = %q", network.RPCEndpoint))
		result = append(result, fmt.Sprintf("chain_id = %d", network.ChainID))
		result = append(result, fmt.Sprintf("symbol = %q", network.Symbol))
		result = append(result, fmt.Sprintf("explorer = %q", network.Explorer))
		result = append(result, fmt.Sprintf("is_active = %t", network.IsActive))
	}

	return result
}

// SanitizeNetworkKey sanitiza uma chave para garantir que seja válida para TOML
func (tsm *TOMLSectionManager) SanitizeNetworkKey(key string) string {
	// Substitui caracteres inválidos por underscore
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, key)
}

// GenerateNetworkKey gera uma chave única e válida para uma rede
func (tsm *TOMLSectionManager) GenerateNetworkKey(name string, chainID int64) string {
	// Sanitiza o nome
	sanitizedName := tsm.SanitizeNetworkKey(name)

	// Cria a chave no formato custom_{sanitized_name}_{chain_id}
	return fmt.Sprintf("custom_%s_%d", sanitizedName, chainID)
}
