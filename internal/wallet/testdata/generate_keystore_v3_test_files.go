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
	"strings"

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

	// Generate additional valid keystores with different parameters
	testKeystores = append(testKeystores, generateAdditionalValidKeystores()...)

	// Generate additional invalid keystores for negative testing
	testKeystores = append(testKeystores, generateAdditionalInvalidKeystores()...)

	// Generate documentation
	updateDocumentation(testKeystores)

	fmt.Println("All additional test keystore files generated successfully!")
}

func generateAdditionalValidKeystores() []TestKeystoreInfo {
	var testKeystores []TestKeystoreInfo

	// Generate keystore with very strong scrypt parameters
	testKeystores = append(testKeystores, generateRealKeystore(
		"real_keystore_v3_strong_scrypt.json",
		"testpassword",
		"8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e",
		keystore.StandardScryptN*2, // Double the standard parameters
		keystore.StandardScryptP*2,
		"Strong scrypt parameters (more secure but slower)",
	))

	// Generate keystore with special characters in the password
	testKeystores = append(testKeystores, generateRealKeystore(
		"real_keystore_v3_special_chars_password.json",
		"!@#$%^&*()_+{}|:<>?~",
		"9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f",
		keystore.StandardScryptN,
		keystore.StandardScryptP,
		"Password with only special characters",
	))

	// Generate keystore with very long password
	longPassword := strings.Repeat("LongPassword123!@#", 10) // 160 characters
	testKeystores = append(testKeystores, generateRealKeystore(
		"real_keystore_v3_long_password.json",
		longPassword,
		"0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a",
		keystore.StandardScryptN,
		keystore.StandardScryptP,
		"Very long password (160 characters)",
	))

	// Generate keystore with Unicode password
	testKeystores = append(testKeystores, generateRealKeystore(
		"real_keystore_v3_unicode_password.json",
		"пароль密码パスワードكلمة السر", // Password in Russian, Chinese, Japanese, Arabic
		"1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b",
		keystore.StandardScryptN,
		keystore.StandardScryptP,
		"Password with Unicode characters",
	))

	// Generate keystore with PBKDF2 and high iteration count
	testKeystores = append(testKeystores, generatePBKDF2Keystore(
		"real_keystore_v3_pbkdf2_high_iterations.json",
		"testpassword",
		"2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c",
		"PBKDF2 with high iteration count (20480)",
		20480, // Double the standard iterations
	))

	return testKeystores
}

func generateAdditionalInvalidKeystores() []TestKeystoreInfo {
	var testKeystores []TestKeystoreInfo

	// Generate keystore with non-standard cipher
	testKeystores = append(testKeystores, generateNonStandardCipher())

	// Generate keystore with invalid address checksum
	testKeystores = append(testKeystores, generateInvalidAddressChecksum())

	// Generate keystore with string version instead of number
	testKeystores = append(testKeystores, generateStringVersion())

	// Generate keystore with missing KDF
	testKeystores = append(testKeystores, generateMissingKDF())

	// Skip the following for now as they have issues
	// Generate keystore with floating point version number
	// testKeystores = append(testKeystores, generateFloatVersion())

	// Generate keystore with extremely large scrypt parameters
	// testKeystores = append(testKeystores, generateExtremeScryptParams())

	// Generate keystore with floating point version number
	testKeystores = append(testKeystores, generateFloatVersion())

	// Generate keystore with extremely large scrypt parameters
	testKeystores = append(testKeystores, generateExtremeScryptParams())

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

func generatePBKDF2Keystore(filename, password, privateKeyHex, notes string, iterations int) TestKeystoreInfo {
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

	// Use the go-ethereum keystore library to encrypt with scrypt first
	keystoreJSON, err := keystore.EncryptKey(key, password, keystore.LightScryptN, keystore.LightScryptP)
	if err != nil {
		log.Fatalf("Failed to encrypt key: %v", err)
	}

	// Parse the keystore JSON
	var keystoreMap map[string]interface{}
	json.Unmarshal(keystoreJSON, &keystoreMap)

	// Modify to use PBKDF2
	crypto := keystoreMap["crypto"].(map[string]interface{})
	crypto["kdf"] = "pbkdf2"

	// Create PBKDF2 parameters
	kdfParams := make(map[string]interface{})
	kdfParams["c"] = iterations // Custom iteration count
	kdfParams["dklen"] = 32     // Derived key length
	kdfParams["prf"] = "hmac-sha256"
	kdfParams["salt"] = generateRandomHex(32)

	crypto["kdfparams"] = kdfParams

	// Write back to file
	modifiedJSON, _ := json.MarshalIndent(keystoreMap, "", "  ")
	outputPath := filepath.Join("keystores", filename)
	ioutil.WriteFile(outputPath, modifiedJSON, 0644)

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

func generateNonStandardCipher() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d"
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to convert to ECDSA: %v", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressHex := address.Hex()

	// Create a temporary keystore
	tempDir, err := ioutil.TempDir("", "keystore")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ks := keystore.NewKeyStore(tempDir, keystore.LightScryptN, keystore.LightScryptP)
	account, err := ks.ImportECDSA(privateKey, "testpassword")
	if err != nil {
		log.Fatalf("Failed to import private key: %v", err)
	}

	keystoreJSON, err := ioutil.ReadFile(account.URL.Path)
	if err != nil {
		log.Fatalf("Failed to read keystore file: %v", err)
	}

	// Parse the keystore JSON
	var keystoreMap map[string]interface{}
	err = json.Unmarshal(keystoreJSON, &keystoreMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal keystore JSON: %v", err)
	}

	// Change the cipher to a non-standard one
	crypto := keystoreMap["crypto"].(map[string]interface{})
	crypto["cipher"] = "aes-256-gcm" // Not the standard aes-128-ctr

	// Write back to file
	modifiedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal modified keystore: %v", err)
	}

	filename := "non_standard_cipher.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write modified keystore file: %v", err)
	}

	fmt.Printf("Created non-standard cipher keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   "testpassword",
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        "scrypt",
		Valid:      false,
		Notes:      "Keystore with non-standard cipher algorithm (aes-256-gcm)",
	}
}

func generateInvalidAddressChecksum() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5"
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to convert to ECDSA: %v", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressHex := address.Hex()

	// Create a temporary keystore
	tempDir, err := ioutil.TempDir("", "keystore")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ks := keystore.NewKeyStore(tempDir, keystore.LightScryptN, keystore.LightScryptP)
	account, err := ks.ImportECDSA(privateKey, "testpassword")
	if err != nil {
		log.Fatalf("Failed to import private key: %v", err)
	}

	keystoreJSON, err := ioutil.ReadFile(account.URL.Path)
	if err != nil {
		log.Fatalf("Failed to read keystore file: %v", err)
	}

	// Parse the keystore JSON
	var keystoreMap map[string]interface{}
	err = json.Unmarshal(keystoreJSON, &keystoreMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal keystore JSON: %v", err)
	}

	// Change the address to have an invalid checksum (uppercase some letters)
	// First get the address without 0x prefix
	addrStr := strings.TrimPrefix(addressHex, "0x")

	// Create an address with mixed case (invalid checksum)
	invalidAddr := ""
	for i, c := range addrStr {
		if i%2 == 0 {
			invalidAddr += strings.ToUpper(string(c))
		} else {
			invalidAddr += string(c)
		}
	}

	keystoreMap["address"] = invalidAddr

	// Write back to file
	modifiedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal modified keystore: %v", err)
	}

	filename := "invalid_address_checksum.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write modified keystore file: %v", err)
	}

	fmt.Printf("Created invalid address checksum keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   "testpassword",
		PrivateKey: privateKeyHex,
		Address:    addressHex, // The correct address
		KDF:        "scrypt",
		Valid:      false,
		Notes:      "Keystore with invalid address checksum (mixed case)",
	}
}

func generateStringVersion() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6"
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to convert to ECDSA: %v", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressHex := address.Hex()

	// Create a temporary keystore
	tempDir, err := ioutil.TempDir("", "keystore")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ks := keystore.NewKeyStore(tempDir, keystore.LightScryptN, keystore.LightScryptP)
	account, err := ks.ImportECDSA(privateKey, "testpassword")
	if err != nil {
		log.Fatalf("Failed to import private key: %v", err)
	}

	keystoreJSON, err := ioutil.ReadFile(account.URL.Path)
	if err != nil {
		log.Fatalf("Failed to read keystore file: %v", err)
	}

	// Parse the keystore JSON
	var keystoreMap map[string]interface{}
	err = json.Unmarshal(keystoreJSON, &keystoreMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal keystore JSON: %v", err)
	}

	// Change the version to a string instead of a number
	keystoreMap["version"] = "3" // Should be a number, not a string

	// Write back to file
	modifiedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal modified keystore: %v", err)
	}

	filename := "string_version.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write modified keystore file: %v", err)
	}

	fmt.Printf("Created string version keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   "testpassword",
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        "scrypt",
		Valid:      false,
		Notes:      "Keystore with version as a string instead of a number",
	}
}

func generateFloatVersion() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7"
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to convert to ECDSA: %v", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressHex := address.Hex()

	// Create a temporary keystore
	tempDir, err := ioutil.TempDir("", "keystore")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ks := keystore.NewKeyStore(tempDir, keystore.LightScryptN, keystore.LightScryptP)
	account, err := ks.ImportECDSA(privateKey, "testpassword")
	if err != nil {
		log.Fatalf("Failed to import private key: %v", err)
	}

	keystoreJSON, err := ioutil.ReadFile(account.URL.Path)
	if err != nil {
		log.Fatalf("Failed to read keystore file: %v", err)
	}

	// Parse the keystore JSON
	var keystoreMap map[string]interface{}
	err = json.Unmarshal(keystoreJSON, &keystoreMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal keystore JSON: %v", err)
	}

	// Change the version to a float instead of an integer
	keystoreMap["version"] = 3.0 // Should be an integer, not a float

	// Write back to file
	modifiedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal modified keystore: %v", err)
	}

	filename := "float_version.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write modified keystore file: %v", err)
	}

	fmt.Printf("Created float version keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   "testpassword",
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        "scrypt",
		Valid:      false,
		Notes:      "Keystore with version as a float instead of an integer",
	}
}

func generateMissingKDF() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8"
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to convert to ECDSA: %v", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressHex := address.Hex()

	// Create a temporary keystore
	tempDir, err := ioutil.TempDir("", "keystore")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ks := keystore.NewKeyStore(tempDir, keystore.LightScryptN, keystore.LightScryptP)
	account, err := ks.ImportECDSA(privateKey, "testpassword")
	if err != nil {
		log.Fatalf("Failed to import private key: %v", err)
	}

	keystoreJSON, err := ioutil.ReadFile(account.URL.Path)
	if err != nil {
		log.Fatalf("Failed to read keystore file: %v", err)
	}

	// Parse the keystore JSON
	var keystoreMap map[string]interface{}
	err = json.Unmarshal(keystoreJSON, &keystoreMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal keystore JSON: %v", err)
	}

	// Remove the KDF field
	crypto := keystoreMap["crypto"].(map[string]interface{})
	delete(crypto, "kdf")

	// Write back to file
	modifiedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal modified keystore: %v", err)
	}

	filename := "missing_kdf.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write modified keystore file: %v", err)
	}

	fmt.Printf("Created missing KDF keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   "testpassword",
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        "",
		Valid:      false,
		Notes:      "Keystore missing the crypto.kdf field",
	}
}

func generateExtremeScryptParams() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9"
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to convert to ECDSA: %v", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressHex := address.Hex()

	// Create a temporary keystore
	tempDir, err := ioutil.TempDir("", "keystore")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ks := keystore.NewKeyStore(tempDir, keystore.LightScryptN, keystore.LightScryptP)
	account, err := ks.ImportECDSA(privateKey, "testpassword")
	if err != nil {
		log.Fatalf("Failed to import private key: %v", err)
	}

	keystoreJSON, err := ioutil.ReadFile(account.URL.Path)
	if err != nil {
		log.Fatalf("Failed to read keystore file: %v", err)
	}

	// Parse the keystore JSON
	var keystoreMap map[string]interface{}
	err = json.Unmarshal(keystoreJSON, &keystoreMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal keystore JSON: %v", err)
	}

	// Set extreme scrypt parameters
	crypto := keystoreMap["crypto"].(map[string]interface{})
	kdfparams := crypto["kdfparams"].(map[string]interface{})
	kdfparams["n"] = 16777216 // Extremely high value, would cause memory issues
	kdfparams["p"] = 64       // Extremely high value, would cause CPU issues

	// Write back to file
	modifiedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal modified keystore: %v", err)
	}

	filename := "extreme_scrypt_params.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write modified keystore file: %v", err)
	}

	fmt.Printf("Created extreme scrypt params keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   "testpassword",
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        "scrypt",
		Valid:      false,
		Notes:      "Keystore with extremely high scrypt parameters that could cause resource exhaustion",
	}
}

func generateRandomHex(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func updateDocumentation(newKeystores []TestKeystoreInfo) {
	// Create a new documentation file for the additional keystores
	readmeContent := "# Additional Test Keystore Files\n\n"
	readmeContent += "This file documents additional keystore files generated for testing the keystore validation and import functionality.\n\n"

	// Add valid keystores section
	readmeContent += "## Additional Valid Keystore Files\n\n"
	for _, ks := range newKeystores {
		if ks.Valid {
			readmeContent += fmt.Sprintf("- `%s`: %s\n", ks.Filename, ks.Notes)
			readmeContent += fmt.Sprintf("  - Password: `%s`\n", ks.Password)
			readmeContent += fmt.Sprintf("  - Address: `%s`\n", ks.Address)
			readmeContent += fmt.Sprintf("  - KDF: %s\n\n", ks.KDF)
		}
	}

	// Add invalid keystores section
	readmeContent += "## Additional Invalid Keystore Files\n\n"
	for _, ks := range newKeystores {
		if !ks.Valid {
			readmeContent += fmt.Sprintf("- `%s`: %s\n", ks.Filename, ks.Notes)
			if ks.Password != "" {
				readmeContent += fmt.Sprintf("  - Password: `%s` (for reference only, file is invalid)\n", ks.Password)
			}
			if ks.Address != "" {
				readmeContent += fmt.Sprintf("  - Address: `%s` (for reference only, file is invalid)\n", ks.Address)
			}
			readmeContent += "\n"
		}
	}

	// Write the documentation file
	outputPath := filepath.Join("keystores", "ADDITIONAL_KEYSTORES.md")
	err := ioutil.WriteFile(outputPath, []byte(readmeContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write documentation file: %v", err)
	}

	fmt.Println("Created documentation file for additional keystores")
}
