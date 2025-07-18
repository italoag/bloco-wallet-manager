# Test Keystore Files

This directory contains sample keystore files for testing the keystore validation and import functionality.

## Valid Keystore Files

| Filename | Password | Address | KDF | Notes |
|----------|----------|---------|-----|-------|
| `real_keystore_v3_standard.json` | `testpassword` | `0xAf6D46d1E55AA87772Fb1538FE4d36AAA70f4e06` | scrypt | Standard scrypt parameters |
| `real_keystore_v3_light.json` | `testpassword` | `0x44BD130B9F2032705e2B3C84b01e1305941c6312` | scrypt | Light scrypt parameters (faster but less secure) |
| `real_keystore_v3_complex_password.json` | `P@$w0rd!123#ComplexPassword` | `0xF32f7C95CD7f616674Cb06d5E253CAC345E2722B` | scrypt | Complex password with special characters |
| `real_keystore_v3_empty_password.json` | `` | `0xBE9392958b1d9a6145f4Ed4531f0055863C3ecd8` | scrypt | Empty password (not recommended but should work) |
| `real_keystore_v3_pbkdf2.json` | `testpassword` | `0xF3a434F00C66A6827ba72a12fCA3fA7c219E1692` | pbkdf2 | PBKDF2 key derivation function |

## Invalid Keystore Files

| Filename | Issue | Notes |
|----------|-------|-------|
| `corrupted_ciphertext.json` | Invalid | Corrupted ciphertext that is not valid hex |
| `invalid_mac.json` | Invalid | Invalid MAC value (will fail password validation) |
| `invalid_kdf_params.json` | Invalid | Invalid KDF parameters (negative scrypt N value) |
| `invalid_version.json` | Invalid | Keystore with version 2 instead of 3 |
| `invalid_json.json` | Invalid | Keystore with invalid JSON syntax |
| `missing_address.json` | Invalid | Keystore missing the address field |
| `invalid_address.json` | Invalid | Keystore with an invalid address format |
| `missing_crypto.json` | Invalid | Keystore missing the crypto field |
| `missing_cipher.json` | Invalid | Keystore missing the crypto.cipher field |
| `missing_ciphertext.json` | Invalid | Keystore missing the crypto.ciphertext field |
| `missing_iv.json` | Invalid | Keystore missing the crypto.cipherparams.iv field |
| `missing_mac.json` | Invalid | Keystore missing the crypto.mac field |
| `missing_scrypt_dklen.json` | Invalid | Keystore missing the crypto.kdfparams.dklen field |
| `unsupported_kdf.json` | Invalid | Keystore with an unsupported KDF algorithm |
| `malformed_address_with_prefix.json` | Invalid | Keystore with a malformed address (0x prefix but incorrect length) |

## Generating Test Keystore Files

The `generate_comprehensive_keystores.go` file can be used to generate test keystore files with known private keys and passwords. This is useful for testing the keystore validation and import functionality.

To generate new keystore files:

```bash
cd internal/wallet/testdata
go run generate_comprehensive_keystores.go
```

## Usage in Tests

These files can be used for testing the keystore validation and import functionality. The valid keystore files can be used to test successful imports, while the invalid keystore files can be used to test error handling.

Example:

```go
func TestImportWalletFromKeystoreV3(t *testing.T) {
    // Test with valid keystore file
    walletDetails, err := walletService.ImportWalletFromKeystoreV3("Test Wallet", "testdata/keystores/real_keystore_v3_standard.json", "testpassword")
    assert.NoError(t, err)
    assert.NotNil(t, walletDetails)

    // Test with invalid keystore file
    walletDetails, err = walletService.ImportWalletFromKeystoreV3("Test Wallet", "testdata/keystores/invalid_version.json", "testpassword")
    assert.Error(t, err)
    assert.Nil(t, walletDetails)
}
```

## Private Keys and Security

For security reasons, the private keys used in these test files should never be used in a production environment. They are included here only for testing purposes and are hardcoded in the generator script.

**Test Private Keys (for reference only):**
- `real_keystore_v3_standard.json`: `1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b`
- `real_keystore_v3_light.json`: `2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c`
- `real_keystore_v3_complex_password.json`: `3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d`
- `real_keystore_v3_empty_password.json`: `4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e`
- `real_keystore_v3_pbkdf2.json`: `5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f`

**IMPORTANT:** Never use these private keys in production or with real funds!
