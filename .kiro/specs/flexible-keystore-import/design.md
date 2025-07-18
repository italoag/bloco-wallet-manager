# Design Document

## Overview

Este documento descreve o design para implementar a funcionalidade de importação flexível de arquivos keystore v3 no BlocoWallet, permitindo a importação de arquivos com qualquer extensão ou sem extensão. O design foca em modificar a validação de arquivos para verificar o conteúdo JSON em vez de depender da extensão do arquivo.

## Architecture

### Current Architecture Issues
- Validação rígida que exige extensão .json para arquivos keystore
- Rejeição de arquivos válidos baseada apenas na extensão do arquivo
- Mensagens de erro que enfatizam a extensão do arquivo em vez do conteúdo

### Proposed Architecture
- Validação baseada no conteúdo do arquivo em vez da extensão
- Verificação de JSON válido independentemente da extensão
- Mensagens de erro mais claras sobre o formato do conteúdo

## Components and Interfaces

### 1. Modificação na Interface de Usuário (TUI)

A principal modificação será na função de importação de keystore na interface do usuário, removendo a verificação de extensão .json:

```go
// Antes
if !strings.HasSuffix(strings.ToLower(keystorePath), ".json") {
    // Use specific error type for invalid file type
    keystoreErr := wallet.NewKeystoreImportError(
        wallet.ErrorInvalidKeystore,
        fmt.Sprintf("File %s is not a JSON file", keystorePath),
        nil,
    )
    m.err = errors.Wrap(fmt.Errorf(localization.FormatKeystoreErrorWithField(
        keystoreErr.GetLocalizedMessage(),
        "",
    )), 0)
    return m, nil
}

// Depois
// Remover completamente esta verificação, permitindo qualquer extensão de arquivo
```

### 2. Aprimoramento da Validação de Conteúdo JSON

Reforçar a validação do conteúdo JSON no `KeystoreValidator`:

```go
// ValidateKeystoreV3 parses JSON data and validates the keystore structure
func (kv *KeystoreValidator) ValidateKeystoreV3(data []byte) (*KeystoreV3, error) {
    var keystore KeystoreV3

    // Parse JSON
    if err := json.Unmarshal(data, &keystore); err != nil {
        return nil, NewKeystoreImportError(ErrorInvalidJSON, "O arquivo não contém um JSON válido", err)
    }

    // Validate structure
    if err := kv.ValidateStructure(&keystore); err != nil {
        return nil, err
    }

    return &keystore, nil
}
```

### 3. Atualização das Mensagens de Erro

Modificar as mensagens de erro para focar no conteúdo em vez da extensão:

```go
// Atualizar mensagens de erro no pacote de localização
"keystore_invalid_json": "O arquivo não contém um JSON válido",
"keystore_invalid_structure": "O arquivo não contém um keystore v3 válido",
```

## Data Models

Não há necessidade de alterações nos modelos de dados existentes. A estrutura `KeystoreV3` e tipos relacionados permanecerão os mesmos.

## Error Handling

### Mensagens de Erro Atualizadas

1. **Erro de JSON Inválido**: Focar no conteúdo do arquivo em vez da extensão
   - Antes: "File %s is not a JSON file"
   - Depois: "O arquivo não contém um JSON válido"

2. **Erro de Estrutura Inválida**: Esclarecer que o problema está no formato do keystore
   - Antes: "Invalid keystore structure"
   - Depois: "O arquivo não contém um keystore v3 válido"

### Fluxo de Tratamento de Erros

1. Verificar existência do arquivo
2. Tentar ler o conteúdo do arquivo
3. Validar se o conteúdo é JSON válido
4. Validar se o JSON tem a estrutura de um keystore v3
5. Prosseguir com a validação de campos específicos do keystore

## Testing Strategy

### Unit Tests

#### Testes de Validação de Keystore
- Testar com arquivos keystore v3 válidos com diferentes extensões:
  - `.json` (padrão)
  - `.key`
  - `.keystoremaster`
  - Sem extensão
  - Extensões compostas (como `.jsonmaster`)
  - Nomes de arquivo complexos (como o exemplo fornecido: `3cc7dc4096856c6e8fa5a179ff6acf7cdbb727720x3cc7dc4096856c6e8fa5a179ff6acf7cdbb72772.keystoremaster.jsonmaster.jsonkey`)

- Testar com arquivos inválidos:
  - Arquivos que não são JSON válido
  - Arquivos JSON que não são keystores v3
  - Arquivos vazios

#### Testes de Importação
- Testar o fluxo completo de importação com diferentes tipos de nomes de arquivo
- Verificar se o arquivo é copiado corretamente para o diretório de destino
- Verificar se o mnemônico é gerado corretamente

### Integration Tests

- Testar o fluxo completo de importação através da interface do usuário
- Verificar se a carteira é criada corretamente no banco de dados
- Verificar se o arquivo keystore é copiado para o diretório gerenciado

### Test Data

Criar arquivos de teste adicionais com diferentes extensões:
- `valid_keystore_v3.key`
- `valid_keystore_v3` (sem extensão)
- `valid_keystore_v3.keystoremaster`
- `3cc7dc4096856c6e8fa5a179ff6acf7cdbb727720x3cc7dc4096856c6e8fa5a179ff6acf7cdbb72772.keystoremaster.jsonmaster.jsonkey`

## Implementation Approach

### Fase 1: Remover Validação de Extensão
1. Modificar a função de importação na interface do usuário para remover a verificação de extensão .json
2. Atualizar as mensagens de erro para focar no conteúdo do arquivo em vez da extensão

### Fase 2: Aprimorar Validação de Conteúdo
1. Reforçar a validação de conteúdo JSON no `KeystoreValidator`
2. Garantir que as mensagens de erro sejam claras sobre o problema específico encontrado

### Fase 3: Testes e Validação
1. Criar arquivos de teste com diferentes extensões
2. Implementar testes unitários para validar a importação com diferentes tipos de arquivo
3. Implementar testes de integração para o fluxo completo

## Security Considerations

### Validação de Conteúdo
- Garantir que a validação de conteúdo JSON seja robusta para evitar injeção de código
- Manter os limites de tamanho de arquivo para evitar ataques de negação de serviço

### Manipulação de Arquivos
- Manter as verificações de permissão de arquivo
- Continuar validando os caminhos de arquivo para evitar ataques de travessia de diretório

## Performance Considerations

### Leitura de Arquivos
- Manter o streaming de arquivos grandes em vez de carregá-los inteiramente na memória
- Implementar timeout para operações de arquivo

### Validação de JSON
- Otimizar a validação de JSON para arquivos grandes
- Considerar limites de tamanho para evitar problemas de desempenho

## Backward Compatibility

### Compatibilidade com Keystores Existentes
- Garantir que keystores existentes com extensão .json continuem funcionando
- Manter o formato de saída dos arquivos keystore com extensão .json

### Compatibilidade com APIs
- Garantir que as APIs existentes continuem funcionando sem modificações
- Manter a compatibilidade com ferramentas externas que possam depender do formato atual