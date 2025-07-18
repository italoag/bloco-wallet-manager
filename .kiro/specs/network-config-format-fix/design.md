# Design Document

## Overview

Este documento descreve o design para corrigir os problemas de formatação na configuração de redes no BlocoWallet. O foco é garantir que o arquivo de configuração TOML seja gerado corretamente, sem duplicação de seções e com a formatação adequada para as redes blockchain.

## Architecture

### Current Architecture Issues
- Duplicação da seção [networks] no arquivo de configuração
- Formatação incorreta das subseções de rede
- Falta de atomicidade nas operações de escrita do arquivo
- Ausência de backup do arquivo de configuração

### Proposed Architecture
- Implementação robusta para manipulação de arquivos TOML
- Estratégia de escrita atômica para evitar corrupção de dados
- Validação de chaves e valores antes da escrita
- Sistema de backup e recuperação para operações críticas

## Components and Interfaces

### 1. Config File Manager

```go
// ConfigFileManager gerencia operações de leitura e escrita no arquivo de configuração
type ConfigFileManager struct {
    configPath string
}

// NewConfigFileManager cria uma nova instância do gerenciador
func NewConfigFileManager(configPath string) *ConfigFileManager

// ReadConfig lê o arquivo de configuração atual
func (cfm *ConfigFileManager) ReadConfig() ([]string, error)

// WriteConfig escreve o conteúdo atualizado no arquivo de configuração
func (cfm *ConfigFileManager) WriteConfig(content []string) error

// BackupConfig cria um backup do arquivo de configuração
func (cfm *ConfigFileManager) BackupConfig() (string, error)

// RestoreConfig restaura o arquivo de configuração a partir de um backup
func (cfm *ConfigFileManager) RestoreConfig(backupPath string) error
```

### 2. TOML Section Manager

```go
// TOMLSectionManager gerencia seções em arquivos TOML
type TOMLSectionManager struct{}

// FindSection encontra o início e fim de uma seção específica
func (tsm *TOMLSectionManager) FindSection(lines []string, sectionName string) (start, end int)

// RemoveSection remove uma seção específica do conteúdo
func (tsm *TOMLSectionManager) RemoveSection(lines []string, sectionName string) []string

// AddSection adiciona uma nova seção ao conteúdo
func (tsm *TOMLSectionManager) AddSection(lines []string, sectionName string, content []string) []string

// FormatNetworkSection formata a seção de redes corretamente
func (tsm *TOMLSectionManager) FormatNetworkSection(networks map[string]config.Network) []string
```

### 3. Enhanced Network Configuration Manager

```go
// NetworkConfigManager gerencia a configuração de redes
type NetworkConfigManager struct {
    configManager *ConfigFileManager
    sectionManager *TOMLSectionManager
}

// SaveNetworksConfig salva a configuração de redes no arquivo
func (ncm *NetworkConfigManager) SaveNetworksConfig(cfg *config.Config) error

// GenerateNetworkKey gera uma chave única e válida para uma rede
func (ncm *NetworkConfigManager) GenerateNetworkKey(name string, chainID int64) string

// ValidateNetworkKey valida se uma chave de rede é válida para TOML
func (ncm *NetworkConfigManager) ValidateNetworkKey(key string) bool

// SanitizeNetworkKey sanitiza uma chave para garantir que seja válida
func (ncm *NetworkConfigManager) SanitizeNetworkKey(key string) string
```

## Data Models

### Network Configuration Format
- **[networks]**: Seção principal para configurações de rede
- **[networks.{key}]**: Subseção para cada rede individual
  - **name**: Nome da rede (string)
  - **rpc_endpoint**: Endpoint RPC da rede (string)
  - **chain_id**: ID da cadeia (int64)
  - **symbol**: Símbolo da moeda (string)
  - **explorer**: URL do explorador de blocos (string)
  - **is_active**: Status de ativação da rede (boolean)

### Network Key Format
- Formato: `custom_{sanitized_name}_{chain_id}`
- Caracteres permitidos: letras, números e underscore
- Caracteres inválidos são substituídos por underscore

## Error Handling

### Error Classification
1. **File System Errors**: Erros de leitura/escrita de arquivo
2. **Format Errors**: Erros de formatação TOML
3. **Validation Errors**: Erros de validação de chaves e valores
4. **Backup Errors**: Erros relacionados a backup e restauração

### Error Recovery
- Backup automático antes de modificações críticas
- Restauração do backup em caso de falha
- Validação do arquivo após escrita
- Mensagens de erro específicas para cada tipo de falha

## Testing Strategy

### Unit Tests
- **Config File Manager Tests**
  - Testes de leitura e escrita de arquivo
  - Testes de backup e restauração
  - Testes de tratamento de erros

- **TOML Section Manager Tests**
  - Testes de manipulação de seções
  - Testes de formatação de seções de rede
  - Testes de validação de chaves

- **Network Configuration Manager Tests**
  - Testes de geração de chaves
  - Testes de sanitização de chaves
  - Testes de salvamento de configuração

### Integration Tests
- **End-to-End Configuration Tests**
  - Testes de adição, edição e remoção de redes
  - Testes de carregamento de configuração
  - Testes de consistência de dados

## Implementation Approach

### Phase 1: Refatoração do Gerenciamento de Arquivos
1. Implementar ConfigFileManager com operações atômicas
2. Adicionar sistema de backup e restauração
3. Melhorar tratamento de erros

### Phase 2: Melhoria da Formatação TOML
1. Implementar TOMLSectionManager
2. Corrigir formatação da seção [networks]
3. Garantir unicidade das seções

### Phase 3: Aprimoramento da Geração de Chaves
1. Implementar geração robusta de chaves
2. Adicionar validação e sanitização de chaves
3. Garantir compatibilidade com formato TOML

### Phase 4: Integração com UI
1. Atualizar funções de UI para usar os novos gerenciadores
2. Melhorar feedback de erro para o usuário
3. Garantir consistência da interface

## Security Considerations

### File Handling
- Operações atômicas para evitar corrupção de dados
- Validação de permissões de arquivo
- Tratamento seguro de caminhos de arquivo

### Data Validation
- Validação de entrada antes da escrita
- Sanitização de chaves e valores
- Prevenção de injeção de caracteres especiais

## Performance Considerations

### File Operations
- Minimizar operações de I/O
- Usar buffers para operações de leitura/escrita
- Implementar cache quando apropriado

### Memory Usage
- Processar arquivos linha por linha quando possível
- Evitar carregar arquivos grandes inteiramente na memória
- Limpar recursos após uso

## Backward Compatibility

### Configuration Format
- Manter compatibilidade com formato atual
- Garantir que configurações existentes sejam preservadas
- Implementar migração automática quando necessário

### API Compatibility
- Manter assinaturas de função existentes
- Adicionar novas funcionalidades de forma não disruptiva
- Documentar mudanças para desenvolvedores