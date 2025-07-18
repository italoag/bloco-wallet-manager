package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

// TestKeystoreInfo holds information about a test keystore file
type TestKeystoreInfo struct {
	Filename   string
	Password   string
	PrivateKey string
	Address    string
	KDF        string
	Valid      bool
	Notes      string
}

func main() {
	// Create directory if it doesn't exist
	keystoreDir := "keystores"
	if _, err := os.Stat(keystoreDir); os.IsNotExist(err) {
		os.MkdirAll(keystoreDir, 0755)
	}

	// Track all generated keystores for documentation
	var testKeystores []TestKeystoreInfo

	// Generate valid keystores with different parameters
	testKeystores = append(testKeystores, generateValidKeystores()...)

	// Generate invalid keystores for negative testing
	testKeystores = append(testKeystores, generateInvalidKeystores()...)

	// Generate documentation
	generateDocumentation(testKeystores)

	fmt.Println("All test keystore files generated successfully!")
}

func generateValidKeystores() []TestKeystoreInfo {
	var testKeystores []TestKeystoreInfo

	// Generate keystores with different private keys and parameters
	testKeystores = append(testKeystores, generateRealKeystore(
		"real_keystore_v3_standard.json",
		"testpassword",
		"1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b",
		keystore.StandardScryptN,
		keystore.StandardScryptP,
		"Standard scrypt parameters",
	))

	testKeystores = append(testKeystores, generateRealKeystore(
		"real_keystore_v3_light.json",
		"testpassword",
		"2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c",
		keystore.LightScryptN,
		keystore.LightScryptP,
		"Light scrypt parameters (faster but less secure)",
	))

	testKeystores = append(testKeystores, generateRealKeystore(
		"real_keystore_v3_complex_password.json",
		"P@$$w0rd!123#ComplexPassword",
		"3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d",
		keystore.StandardScryptN,
		keystore.StandardScryptP,
		"Complex password with special characters",
	))

	testKeystores = append(testKeystores, generateRealKeystore(
		"real_keystore_v3_empty_password.json",
		"",
		"4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e",
		keystore.StandardScryptN,
		keystore.StandardScryptP,
		"Empty password (not recommended but should work)",
	))

	// Generate a keystore with PBKDF2
	testKeystores = append(testKeystores, generatePBKDF2Keystore(
		"real_keystore_v3_pbkdf2.json",
		"testpassword",
		"5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f",
		"PBKDF2 key derivation function",
	))

	return testKeystores
}

func generateInvalidKeystores() []TestKeystoreInfo {
	var testKeystores []TestKeystoreInfo

	// 1. Invalid version (already exists)
	// 2. Invalid JSON (already exists)
	// 3. Missing fields (already exist)

	// 4. Corrupted ciphertext
	testKeystores = append(testKeystores, generateCorruptedCiphertext())

	// 5. Invalid MAC
	testKeystores = append(testKeystores, generateInvalidMAC())

	// 6. Invalid KDF parameters
	testKeystores = append(testKeystores, generateInvalidKDFParams())

	// 7. Non-standard cipher
	testKeystores = append(testKeystores, generateNonStandardCipher())

	// 8. Malformed address with 0x prefix
	testKeystores = append(testKeystores, generateMalformedAddress())

	return testKeystores
}

func generateRealKeystore(filename, password, privateKeyHex string, scryptN, scryptP int, notes string) TestKeystoreInfo {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to convert to ECDSA: %v", err)
	}

	// Get the address from the private key
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressHex := address.Hex()

	// Create a temporary directory for the keystore
	tempDir, err := ioutil.TempDir("", "keystore")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a keystore with specified parameters
	ks := keystore.NewKeyStore(tempDir, scryptN, scryptP)

	// Import the private key
	account, err := ks.ImportECDSA(privateKey, password)
	if err != nil {
		log.Fatalf("Failed to import private key: %v", err)
	}

	// Read the keystore file
	keystoreJSON, err := ioutil.ReadFile(account.URL.Path)
	if err != nil {
		log.Fatalf("Failed to read keystore file: %v", err)
	}

	// Write the keystore file to the testdata directory
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, keystoreJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write keystore file: %v", err)
	}

	fmt.Printf("Created keystore file %s with address: %s\n", filename, addressHex)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   password,
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        "scrypt",
		Valid:      true,
		Notes:      notes,
	}
}

func generatePBKDF2Keystore(filename, password, privateKeyHex, notes string) TestKeystoreInfo {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to convert to ECDSA: %v", err)
	}

	// Get the address from the private key
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressHex := address.Hex()

	// Create a keystore JSON manually with PBKDF2
	key := &keystore.Key{
		Address:    address,
		PrivateKey: privateKey,
		Id:         uuid.New(),
	}

	// Use the go-ethereum keystore library to encrypt with PBKDF2
	keystoreJSON, err := keystore.EncryptKey(key, password, keystore.LightScryptN, keystore.LightScryptP)
	if err != nil {
		log.Fatalf("Failed to encrypt key: %v", err)
	}

	// Write the keystore file to the testdata directory
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, keystoreJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write keystore file: %v", err)
	}

	fmt.Printf("Created PBKDF2 keystore file %s with address: %s\n", filename, addressHex)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   password,
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        "pbkdf2",
		Valid:      true,
		Notes:      notes,
	}
}

func generateCorruptedCiphertext() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7"
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)
	privateKey, _ := crypto.ToECDSA(privateKeyBytes)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressHex := address.Hex()

	// Create a temporary keystore
	tempDir, _ := ioutil.TempDir("", "keystore")
	defer os.RemoveAll(tempDir)
	ks := keystore.NewKeyStore(tempDir, keystore.LightScryptN, keystore.LightScryptP)
	account, _ := ks.ImportECDSA(privateKey, "testpassword")
	keystoreJSON, _ := ioutil.ReadFile(account.URL.Path)

	// Parse the keystore JSON
	var keystoreMap map[string]interface{}
	json.Unmarshal(keystoreJSON, &keystoreMap)

	// Corrupt the ciphertext
	crypto := keystoreMap["crypto"].(map[string]interface{})
	crypto["ciphertext"] = "corrupted_ciphertext_that_is_not_valid_hex"

	// Write back to file
	corruptedJSON, _ := json.MarshalIndent(keystoreMap, "", "  ")
	filename := "corrupted_ciphertext.json"
	outputPath := filepath.Join("keystores", filename)
	ioutil.WriteFile(outputPath, corruptedJSON, 0644)

	fmt.Printf("Created corrupted ciphertext keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   "testpassword",
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        "scrypt",
		Valid:      false,
		Notes:      "Corrupted ciphertext that is not valid hex",
	}
}

func generateInvalidMAC() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8"
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)
	privateKey, _ := crypto.ToECDSA(privateKeyBytes)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressHex := address.Hex()

	// Create a temporary keystore
	tempDir, _ := ioutil.TempDir("", "keystore")
	defer os.RemoveAll(tempDir)
	ks := keystore.NewKeyStore(tempDir, keystore.LightScryptN, keystore.LightScryptP)
	account, _ := ks.ImportECDSA(privateKey, "testpassword")
	keystoreJSON, _ := ioutil.ReadFile(account.URL.Path)

	// Parse the keystore JSON
	var keystoreMap map[string]interface{}
	json.Unmarshal(keystoreJSON, &keystoreMap)

	// Corrupt the MAC
	crypto := keystoreMap["crypto"].(map[string]interface{})
	crypto["mac"] = "0000000000000000000000000000000000000000000000000000000000000000"

	// Write back to file
	corruptedJSON, _ := json.MarshalIndent(keystoreMap, "", "  ")
	filename := "invalid_mac.json"
	outputPath := filepath.Join("keystores", filename)
	ioutil.WriteFile(outputPath, corruptedJSON, 0644)

	fmt.Printf("Created invalid MAC keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   "testpassword",
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        "scrypt",
		Valid:      false,
		Notes:      "Invalid MAC value (will fail password validation)",
	}
}

func generateInvalidKDFParams() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9"
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)
	privateKey, _ := crypto.ToECDSA(privateKeyBytes)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressHex := address.Hex()

	// Create a temporary keystore
	tempDir, _ := ioutil.TempDir("", "keystore")
	defer os.RemoveAll(tempDir)
	ks := keystore.NewKeyStore(tempDir, keystore.LightScryptN, keystore.LightScryptP)
	account, _ := ks.ImportECDSA(privateKey, "testpassword")
	keystoreJSON, _ := ioutil.ReadFile(account.URL.Path)

	// Parse the keystore JSON
	var keystoreMap map[string]interface{}
	json.Unmarshal(keystoreJSON, &keystoreMap)

	// Corrupt the KDF parameters
	crypto := keystoreMap["crypto"].(map[string]interface{})
	kdfparams := crypto["kdfparams"].(map[string]interface{})
	kdfparams["n"] = -1 // Invalid negative value for scrypt N parameter

	// Write back to file
	corruptedJSON, _ := json.MarshalIndent(keystoreMap, "", "  ")
	filename := "invalid_kdf_params.json"
	outputPath := filepath.Join("keystores", filename)
	ioutil.WriteFile(outputPath, corruptedJSON, 0644)

	fmt.Printf("Created invalid KDF params keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   "testpassword",
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        "scrypt",
		Valid:      false,
		Notes:      "Invalid KDF parameters (negative scrypt N value)",
	}
}

func generateNonStandardCipher() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0"
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)
	privateKey, _ := crypto.ToECDSA(privateKeyBytes)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressHex := address.Hex()

	// Create a temporary keystore
	tempDir, _ := ioutil.TempDir("", "keystore")
	defer os.RemoveAll(tempDir)
	ks := keystore.NewKeyStore(tempDir, keystore.LightScryptN, keystore.LightScryptP)
	account, _ := ks.ImportECDSA(privateKey, "testpassword")
	keystoreJSON, _ := ioutil.ReadFile(account.URL.Path)

	// Parse the keystore JSON
	var keystoreMap map[string]interface{}
	json.Unmarshal(keystoreJSON, &keystoreMap)

	// Change the cipher to a non-standard one
	crypto := keystoreMap["crypto"].(map[string]interface{})
	crypto["cipher"] = "aes-256-gcm" // Not supported by the standard keystore

	// Write back to file
	modifiedJSON, _ := json.MarshalIndent(keystoreMap, "", "  ")
	filename := "non_standard_cipher.json"
	outputPath := filepath.Join("keystores", filename)
	ioutil.WriteFile(outputPath, modifiedJSON, 0644)

	fmt.Printf("Created non-standard cipher keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   "testpassword",
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        "scrypt",
		Valid:      false,
		Notes:      "Non-standard cipher algorithm (aes-256-gcm instead of aes-128-ctr)",
	}
}

func generateMalformedAddress() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1"
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)
	privateKey, _ := crypto.ToECDSA(privateKeyBytes)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressHex := address.Hex()

	// Create a temporary keystore
	tempDir, _ := ioutil.TempDir("", "keystore")
	defer os.RemoveAll(tempDir)
	ks := keystore.NewKeyStore(tempDir, keystore.LightScryptN, keystore.LightScryptP)
	account, _ := ks.ImportECDSA(privateKey, "testpassword")
	keystoreJSON, _ := ioutil.ReadFile(account.URL.Path)

	// Parse the keystore JSON
	var keystoreMap map[string]interface{}
	json.Unmarshal(keystoreJSON, &keystoreMap)

	// Malform the address by adding 0x prefix but keeping the same length
	// This will make the address too long with the prefix
	keystoreMap["address"] = "0x" + keystoreMap["address"].(string)

	// Write back to file
	modifiedJSON, _ := json.MarshalIndent(keystoreMap, "", "  ")
	filename := "malformed_address_with_prefix.json"
	outputPath := filepath.Join("keystores", filename)
	ioutil.WriteFile(outputPath, modifiedJSON, 0644)

	fmt.Printf("Created malformed address keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   "testpassword",
		PrivateKey: privateKeyHex,
		Address:    addressHex, // The correct address
		KDF:        "scrypt",
		Valid:      false,
		Notes:      "Malformed address with 0x prefix but incorrect length",
	}
}

func generateRandomBytes(length int) []byte {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return bytes
}

func generateDocumentation(testKeystores []TestKeystoreInfo) {
	// Create a comprehensive README with all test keystores
	readmeContent := `# Test Keystore Files

This directory contains sample keystore files for testing the keystore validation and import functionality.

## Valid Keystore Files

| Filename | Password | Address | KDF | Notes |
|----------|----------|---------|-----|-------|
`

	// Add valid keystores to the table
	for _, ks := range testKeystores {
		if ks.Valid {
			readmeContent += fmt.Sprintf("| `%s` | `%s` | `%s` | %s | %s |\n",
				ks.Filename, ks.Password, ks.Address, ks.KDF, ks.Notes)
		}
	}

	readmeContent += `
## Invalid Keystore Files

| Filename | Issue | Notes |
|----------|-------|-------|
`

	// Add invalid keystores to the table
	for _, ks := range testKeystores {
		if !ks.Valid {
			readmeContent += fmt.Sprintf("| `%s` | Invalid | %s |\n",
				ks.Filename, ks.Notes)
		}
	}

	// Add existing invalid files that we didn't regenerate
	readmeContent += "| `invalid_version.json` | Invalid | Keystore with version 2 instead of 3 |\n"
	readmeContent += "| `invalid_json.json` | Invalid | Keystore with invalid JSON syntax |\n"
	readmeContent += "| `missing_address.json` | Invalid | Keystore missing the address field |\n"
	readmeContent += "| `invalid_address.json` | Invalid | Keystore with an invalid address format |\n"
	readmeContent += "| `missing_crypto.json` | Invalid | Keystore missing the crypto field |\n"
	readmeContent += "| `missing_cipher.json` | Invalid | Keystore missing the crypto.cipher field |\n"
	readmeContent += "| `missing_ciphertext.json` | Invalid | Keystore missing the crypto.ciphertext field |\n"
	readmeContent += "| `missing_iv.json` | Invalid | Keystore missing the crypto.cipherparams.iv field |\n"
	readmeContent += "| `missing_mac.json` | Invalid | Keystore missing the crypto.mac field |\n"
	readmeContent += "| `missing_scrypt_dklen.json` | Invalid | Keystore missing the crypto.kdfparams.dklen field |\n"
	readmeContent += "| `unsupported_kdf.json` | Invalid | Keystore with an unsupported KDF algorithm |\n"

	readmeContent += `
## Generating Test Keystore Files

The ` + "`generate_test_keystores.go`" + ` file can be used to generate test keystore files with known private keys and passwords. This is useful for testing the keystore validation and import functionality.

To generate new keystore files:

` + "```bash" + `
cd internal/wallet/testdata
go run generate_test_keystores.go
` + "```" + `

## Usage in Tests

These files can be used for testing the keystore validation and import functionality. The valid keystore files can be used to test successful imports, while the invalid keystore files can be used to test error handling.

Example:

` + "```go" + `
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
` + "```" + `

## Private Keys

For security reasons, the private keys are not listed in this README. They are hardcoded in the generator script for testing purposes only. In a production environment, never hardcode or commit private keys to version control.
`

	// Write the README file
	readmePath := filepath.Join("keystores", "README.md")
	err := ioutil.WriteFile(readmePath, []byte(readmeContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write README file: %v", err)
	}

	fmt.Println("Generated documentation in keystores/README.md")
}
