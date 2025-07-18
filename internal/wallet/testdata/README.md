# Test Keystore Files

This directory contains sample keystore files for testing the keystore validation and import functionality.

## Valid Keystore Files

- `valid_keystore_v3.json`: A valid keystore v3 file with scrypt KDF
  - Password: `testpassword`
  - Address: `0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d`

- `valid_keystore_v3_pbkdf2.json`: A valid keystore v3 file with PBKDF2 KDF
  - Password: `testpassword`
  - Address: `0x5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d`

- `real_keystore_v3.json`: A real keystore v3 file generated with go-ethereum's keystore library (scrypt KDF)
  - Password: `testpassword`
  - Address: Generated from the private key

- `real_keystore_v3_pbkdf2.json`: A real keystore v3 file generated with go-ethereum's keystore library (PBKDF2 KDF)
  - Password: `testpassword`
  - Address: Generated from the private key

- `real_keystore_v3_standard.json`: A real keystore v3 file with standard scrypt parameters
  - Password: `testpassword`
  - Address: Generated from the private key

- `real_keystore_v3_light.json`: A real keystore v3 file with light scrypt parameters
  - Password: `testpassword`
  - Address: Generated from the private key

- `real_keystore_v3_complex_password.json`: A real keystore v3 file with a complex password
  - Password: `P@$$w0rd!123#ComplexPassword`
  - Address: Generated from the private key

- `real_keystore_v3_empty_password.json`: A real keystore v3 file with an empty password
  - Password: `` (empty string)
  - Address: Generated from the private key

## Invalid Keystore Files

- `invalid_version.json`: A keystore file with version 2 instead of 3
- `invalid_json.json`: A keystore file with invalid JSON syntax
- `missing_address.json`: A keystore file missing the address field
- `invalid_address.json`: A keystore file with an invalid address format
- `missing_crypto.json`: A keystore file missing the crypto field
- `missing_cipher.json`: A keystore file missing the crypto.cipher field
- `missing_ciphertext.json`: A keystore file missing the crypto.ciphertext field
- `missing_iv.json`: A keystore file missing the crypto.cipherparams.iv field
- `missing_mac.json`: A keystore file missing the crypto.mac field
- `missing_scrypt_dklen.json`: A keystore file missing the crypto.kdfparams.dklen field
- `unsupported_kdf.json`: A keystore file with an unsupported KDF algorithm
- `corrupted_ciphertext.json`: A keystore file with corrupted ciphertext that is not valid hex
- `invalid_mac.json`: A keystore file with an invalid MAC value (will fail password validation)
- `invalid_kdf_params.json`: A keystore file with invalid KDF parameters (negative scrypt N value)
- `non_standard_cipher.json`: A keystore file with a non-standard cipher algorithm
- `malformed_address_with_prefix.json`: A keystore file with a malformed address (0x prefix but incorrect length)

## Generating Test Keystore Files

There are two generator scripts in this directory:

1. `generate_real_keystore.go`: Generates real keystore files with known private keys and passwords using the go-ethereum keystore library.
2. `generate_test_keystores.go`: Generates a comprehensive set of test keystore files for both valid and invalid scenarios.

To generate new keystore files:

```bash
cd internal/wallet/testdata
go run generate_test_keystores.go
```

This will create various test keystore files in the `keystores` directory and update the README.md file with details about each file.

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

## Security Note

For security reasons, the private keys used in these test files should never be used in a production environment. They are included here only for testing purposes.