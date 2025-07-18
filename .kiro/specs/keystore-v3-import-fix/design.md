# Design Document

## Overview

Este documento descreve o design para corrigir os bugs na funcionalidade de importação de arquivos keystore v3 no BlocoWallet. O design foca em melhorar a validação, tratamento de erros e confiabilidade da importação.

## Architecture

### Current Architecture Issues
- Validação inadequada de formato keystore v3
- Mensagens de erro genéricas
- Geração não-determinística de mnemônico
- Falta de verificação de integridade

### Proposed Architecture
- Camada de validação estruturada para keystores
- Sistema de erros tipados e específicos
- Geração determinística de mnemônico
- Verificação de integridade end-to-end

## Components and Interfaces

### 1. Keystore Validator

```go
type KeystoreValidator struct{}

type KeystoreValidationError struct {
    Type    string
    Message string
    Field   string
}

func (kv *KeystoreValidator) ValidateKeystoreV3(data []byte) (*KeystoreV3, error)
func (kv *KeystoreValidator) ValidateStructure(keystore *KeystoreV3) error
func (kv *KeystoreValidator) ValidateVersion(version interface{}) error
```

### 2. Enhanced Keystore Structure

```go
type KeystoreV3 struct {
    Version int                    `json:"version"`
    ID      string                 `json:"id"`
    Address string                 `json:"address"`
    Crypto  KeystoreV3Crypto      `json:"crypto"`
}

type KeystoreV3Crypto struct {
    Cipher       string                 `json:"cipher"`
    CipherText   string                 `json:"ciphertext"`
    CipherParams KeystoreV3CipherParams `json:"cipherparams"`
    KDF          string                 `json:"kdf"`
    KDFParams    interface{}            `json:"kdfparams"`
    MAC          string                 `json:"mac"`
}
```

### 3. Error Types

```go
type KeystoreImportError struct {
    Type    KeystoreErrorType
    Message string
    Cause   error
}

type KeystoreErrorType int

const (
    ErrorFileNotFound KeystoreErrorType = iota
    ErrorInvalidJSON
    ErrorInvalidKeystore
    ErrorInvalidVersion
    ErrorIncorrectPassword
    ErrorCorruptedFile
    ErrorAddressMismatch
)
```

### 4. Enhanced Import Service

```go
func (ws *WalletService) ImportWalletFromKeystoreV3(name, keystorePath, password string) (*WalletDetails, error) {
    // 1. Validate file existence
    // 2. Read and parse JSON
    // 3. Validate keystore structure
    // 4. Decrypt and verify
    // 5. Generate deterministic mnemonic
    // 6. Create wallet entry
}
```

## Data Models

### Keystore V3 Structure Validation
- **version**: Must be exactly 3
- **address**: Must be valid Ethereum address (with or without 0x prefix)
- **crypto**: Must contain all required encryption parameters
- **id**: UUID identifier (optional but recommended)

### Mnemonic Generation Strategy
- Use deterministic approach based on private key
- Ensure consistency across imports of same keystore
- Use BIP39 standard for mnemonic generation

## Error Handling

### Error Classification
1. **File System Errors**: File not found, permission denied
2. **Format Errors**: Invalid JSON, missing required fields
3. **Validation Errors**: Invalid version, malformed structure
4. **Cryptographic Errors**: Wrong password, corrupted data
5. **Integrity Errors**: Address mismatch, invalid private key

### Error Messages Localization
- Maintain error messages in localization files
- Provide specific, actionable error messages
- Include suggestions for resolution when possible
- Support multiple languages (English, Portuguese, Spanish)
- Dynamically load messages based on current language setting
- Separate validation feedback from recovery suggestions

### Error Recovery
- Graceful handling of all error types
- No partial state changes on failure
- Proper cleanup of temporary files

## Testing Strategy

### Unit Tests
- **Keystore Validation Tests**
  - Valid keystore v3 files
  - Invalid JSON structures
  - Missing required fields
  - Invalid version numbers
  
- **Import Function Tests**
  - Successful import scenarios
  - Various error conditions
  - Password validation
  - File system edge cases

- **Mnemonic Generation Tests**
  - Deterministic generation
  - Consistency across multiple imports
  - BIP39 compliance

### Integration Tests
- **End-to-End Import Flow**
  - Complete import process
  - Database persistence
  - File management
  - UI interaction

### Test Data
- **Sample Keystore Files**
  - Valid keystore v3 files with different encryption parameters
  - Invalid files for negative testing
  - Corrupted files for error handling tests

## Implementation Approach

### Phase 1: Core Validation
1. Implement KeystoreValidator with structure validation
2. Add proper error types and messages
3. Update localization files with specific error messages

### Phase 2: Enhanced Import Logic
1. Refactor ImportWalletFromKeystoreV3 function
2. Add deterministic mnemonic generation
3. Implement integrity verification

### Phase 3: Testing and Validation
1. Create comprehensive test suite
2. Add sample keystore files for testing
3. Validate error handling scenarios

### Phase 4: UI Improvements
1. Update error message display in UI
2. Improve user feedback during import process
3. Add progress indicators for long operations

## Security Considerations

### Password Handling
- Never store passwords in memory longer than necessary
- Use secure comparison for password validation
- Clear sensitive data from memory after use

### File Handling
- Validate file paths to prevent directory traversal
- Ensure proper file permissions on created keystore copies
- Handle temporary files securely

### Cryptographic Operations
- Use established libraries for keystore decryption
- Validate cryptographic parameters before use
- Ensure proper entropy for any random operations

## Performance Considerations

### File Operations
- Stream large files instead of loading entirely in memory
- Implement timeout for file operations
- Cache validation results when appropriate

### Cryptographic Operations
- Optimize keystore decryption process
- Use appropriate key derivation parameters
- Consider async operations for UI responsiveness

## Backward Compatibility

### Existing Wallets
- Ensure existing wallet imports continue to work
- Maintain compatibility with current database schema
- Preserve existing mnemonic encryption approach

### Configuration
- Maintain existing configuration parameters
- Add new validation settings as optional
- Provide migration path for enhanced features