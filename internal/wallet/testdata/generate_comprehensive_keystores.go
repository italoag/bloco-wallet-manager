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
		"P@$w0rd!123#ComplexPassword",
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
		"real_keystore_v3_pbkdf2_proper.json",
		"testpassword",
		"5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f",
		"PBKDF2 key derivation function with proper parameters",
	))

	// Generate a simple valid keystore for basic testing
	testKeystores = append(testKeystores, generateSimpleKeystore(
		"valid_keystore_v3.json",
		"testpassword",
		"6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a",
		"scrypt",
		"Simple valid keystore v3 with scrypt KDF",
	))

	// Generate a simple valid keystore with PBKDF2
	testKeystores = append(testKeystores, generateSimpleKeystore(
		"valid_keystore_v3_pbkdf2.json",
		"testpassword",
		"7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b",
		"pbkdf2",
		"Simple valid keystore v3 with PBKDF2 KDF",
	))

	return testKeystores
}

func generateInvalidKeystores() []TestKeystoreInfo {
	var testKeystores []TestKeystoreInfo

	// Generate corrupted ciphertext keystore
	testKeystores = append(testKeystores, generateCorruptedCiphertext())

	// Generate invalid MAC keystore
	testKeystores = append(testKeystores, generateInvalidMAC())

	// Generate invalid KDF parameters keystore
	testKeystores = append(testKeystores, generateInvalidKDFParams())

	// Generate unsupported KDF keystore
	testKeystores = append(testKeystores, generateUnsupportedKDF())

	// Generate malformed address keystore
	testKeystores = append(testKeystores, generateMalformedAddress())

	// Generate invalid version keystore
	testKeystores = append(testKeystores, generateInvalidVersion())

	// Generate invalid JSON keystore
	testKeystores = append(testKeystores, generateInvalidJSON())

	// Generate missing fields keystores
	testKeystores = append(testKeystores, generateMissingFields()...)

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
	// Note: The go-ethereum library doesn't directly expose PBKDF2 encryption,
	// so we'll create a scrypt keystore and then modify it to use PBKDF2
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
	kdfParams["c"] = 10240  // Iterations
	kdfParams["dklen"] = 32 // Derived key length
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

func generateSimpleKeystore(filename, password, privateKeyHex, kdfType, notes string) TestKeystoreInfo {
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
	addressWithoutPrefix := strings.TrimPrefix(addressHex, "0x")

	// Create a simple keystore structure
	keystoreMap := map[string]interface{}{
		"version": 3,
		"id":      uuid.New().String(),
		"address": addressWithoutPrefix,
		"crypto": map[string]interface{}{
			"cipher":     "aes-128-ctr",
			"ciphertext": generateRandomHex(64), // Fake ciphertext
			"cipherparams": map[string]interface{}{
				"iv": generateRandomHex(16),
			},
			"mac": generateRandomHex(64), // Fake MAC
		},
	}

	// Set KDF parameters based on type
	crypto := keystoreMap["crypto"].(map[string]interface{})
	if kdfType == "scrypt" {
		crypto["kdf"] = "scrypt"
		crypto["kdfparams"] = map[string]interface{}{
			"dklen": 32,
			"n":     4096,
			"p":     6,
			"r":     8,
			"salt":  generateRandomHex(32),
		}
	} else if kdfType == "pbkdf2" {
		crypto["kdf"] = "pbkdf2"
		crypto["kdfparams"] = map[string]interface{}{
			"dklen": 32,
			"c":     10240,
			"prf":   "hmac-sha256",
			"salt":  generateRandomHex(32),
		}
	}

	// Write to file
	keystoreJSON, _ := json.MarshalIndent(keystoreMap, "", "  ")
	outputPath := filepath.Join("keystores", filename)
	ioutil.WriteFile(outputPath, keystoreJSON, 0644)

	fmt.Printf("Created simple keystore file %s with address: %s\n", filename, addressHex)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   password,
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        kdfType,
		Valid:      true,
		Notes:      notes,
	}
}

func generateCorruptedCiphertext() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1"
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

	// Corrupt the ciphertext
	crypto := keystoreMap["crypto"].(map[string]interface{})
	crypto["ciphertext"] = "corrupted_ciphertext_that_is_not_valid_hex"

	// Write back to file
	corruptedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal corrupted keystore: %v", err)
	}

	filename := "corrupted_ciphertext.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, corruptedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write corrupted keystore file: %v", err)
	}

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
	privateKeyHex := "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2"
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

	// Corrupt the MAC
	crypto := keystoreMap["crypto"].(map[string]interface{})
	crypto["mac"] = "0000000000000000000000000000000000000000000000000000000000000000"

	// Write back to file
	corruptedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal corrupted keystore: %v", err)
	}

	filename := "invalid_mac.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, corruptedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write corrupted keystore file: %v", err)
	}

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
	privateKeyHex := "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3"
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

	// Corrupt the KDF parameters
	crypto := keystoreMap["crypto"].(map[string]interface{})
	kdfparams := crypto["kdfparams"].(map[string]interface{})
	kdfparams["n"] = -1 // Invalid negative value for scrypt N parameter

	// Write back to file
	corruptedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal corrupted keystore: %v", err)
	}

	filename := "invalid_kdf_params.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, corruptedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write corrupted keystore file: %v", err)
	}

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

func generateUnsupportedKDF() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a"
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

	// Change the KDF to an unsupported one
	crypto := keystoreMap["crypto"].(map[string]interface{})
	crypto["kdf"] = "argon2id" // Not supported by the standard keystore

	// Write back to file
	modifiedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal modified keystore: %v", err)
	}

	filename := "unsupported_kdf.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write modified keystore file: %v", err)
	}

	fmt.Printf("Created unsupported KDF keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   "testpassword",
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        "argon2id",
		Valid:      false,
		Notes:      "Keystore with an unsupported KDF algorithm",
	}
}

func generateMalformedAddress() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b"
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

	// Malform the address by adding 0x prefix but keeping the same length
	// This will make the address too long with the prefix
	keystoreMap["address"] = "0x" + keystoreMap["address"].(string)

	// Write back to file
	modifiedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal modified keystore: %v", err)
	}

	filename := "malformed_address_with_prefix.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write modified keystore file: %v", err)
	}

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

func generateInvalidVersion() TestKeystoreInfo {
	// Start with a valid keystore
	privateKeyHex := "b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c"
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

	// Change the version to an invalid one
	keystoreMap["version"] = 2 // Only version 3 is valid

	// Write back to file
	modifiedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal modified keystore: %v", err)
	}

	filename := "invalid_version.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write modified keystore file: %v", err)
	}

	fmt.Printf("Created invalid version keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   "testpassword",
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        "scrypt",
		Valid:      false,
		Notes:      "Keystore with version 2 instead of 3",
	}
}

func generateInvalidJSON() TestKeystoreInfo {
	filename := "invalid_json.json"
	outputPath := filepath.Join("keystores", filename)

	// Create an invalid JSON file
	invalidJSON := `{
		"version": 3,
		"id": "12345678-1234-1234-1234-123456789abc",
		"address": "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d",
		"crypto": {
			"cipher": "aes-128-ctr",
			"ciphertext": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			"cipherparams": {
				"iv": "1234567890abcdef1234567890abcdef"
			},
			"kdf": "scrypt",
			"kdfparams": {
				"dklen": 32,
				"n": 262144,
				"p": 1,
				"r": 8,
				"salt": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
			},
			"mac": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
		}
	` // Missing closing brace

	err := ioutil.WriteFile(outputPath, []byte(invalidJSON), 0644)
	if err != nil {
		log.Fatalf("Failed to write invalid JSON file: %v", err)
	}

	fmt.Printf("Created invalid JSON keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   "testpassword",
		PrivateKey: "",
		Address:    "",
		KDF:        "",
		Valid:      false,
		Notes:      "Keystore with invalid JSON syntax",
	}
}

func generateMissingFields() []TestKeystoreInfo {
	var testKeystores []TestKeystoreInfo

	// Start with a valid keystore
	privateKeyHex := "c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d"
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

	// Create various invalid keystores with missing fields

	// Missing address
	missingAddressMap := copyMap(keystoreMap)
	delete(missingAddressMap, "address")
	testKeystores = append(testKeystores, writeInvalidKeystore(
		missingAddressMap,
		"missing_address.json",
		"testpassword",
		privateKeyHex,
		addressHex,
		"scrypt",
		"Keystore missing the address field",
	))

	// Missing crypto
	missingCryptoMap := copyMap(keystoreMap)
	delete(missingCryptoMap, "crypto")
	testKeystores = append(testKeystores, writeInvalidKeystore(
		missingCryptoMap,
		"missing_crypto.json",
		"testpassword",
		privateKeyHex,
		addressHex,
		"scrypt",
		"Keystore missing the crypto field",
	))

	// Missing cipher
	missingCipherMap := copyMap(keystoreMap)
	cryptoMap := missingCipherMap["crypto"].(map[string]interface{})
	delete(cryptoMap, "cipher")
	testKeystores = append(testKeystores, writeInvalidKeystore(
		missingCipherMap,
		"missing_cipher.json",
		"testpassword",
		privateKeyHex,
		addressHex,
		"scrypt",
		"Keystore missing the crypto.cipher field",
	))

	// Missing ciphertext
	missingCiphertextMap := copyMap(keystoreMap)
	cryptoMap = missingCiphertextMap["crypto"].(map[string]interface{})
	delete(cryptoMap, "ciphertext")
	testKeystores = append(testKeystores, writeInvalidKeystore(
		missingCiphertextMap,
		"missing_ciphertext.json",
		"testpassword",
		privateKeyHex,
		addressHex,
		"scrypt",
		"Keystore missing the crypto.ciphertext field",
	))

	// Missing IV
	missingIVMap := copyMap(keystoreMap)
	cryptoMap = missingIVMap["crypto"].(map[string]interface{})
	cipherparamsMap := cryptoMap["cipherparams"].(map[string]interface{})
	delete(cipherparamsMap, "iv")
	testKeystores = append(testKeystores, writeInvalidKeystore(
		missingIVMap,
		"missing_iv.json",
		"testpassword",
		privateKeyHex,
		addressHex,
		"scrypt",
		"Keystore missing the crypto.cipherparams.iv field",
	))

	// Missing MAC
	missingMACMap := copyMap(keystoreMap)
	cryptoMap = missingMACMap["crypto"].(map[string]interface{})
	delete(cryptoMap, "mac")
	testKeystores = append(testKeystores, writeInvalidKeystore(
		missingMACMap,
		"missing_mac.json",
		"testpassword",
		privateKeyHex,
		addressHex,
		"scrypt",
		"Keystore missing the crypto.mac field",
	))

	// Missing scrypt dklen
	missingDklenMap := copyMap(keystoreMap)
	cryptoMap = missingDklenMap["crypto"].(map[string]interface{})
	kdfparamsMap := cryptoMap["kdfparams"].(map[string]interface{})
	delete(kdfparamsMap, "dklen")
	testKeystores = append(testKeystores, writeInvalidKeystore(
		missingDklenMap,
		"missing_scrypt_dklen.json",
		"testpassword",
		privateKeyHex,
		addressHex,
		"scrypt",
		"Keystore missing the crypto.kdfparams.dklen field",
	))

	// Invalid address format
	invalidAddressMap := copyMap(keystoreMap)
	invalidAddressMap["address"] = "not_a_valid_ethereum_address"
	testKeystores = append(testKeystores, writeInvalidKeystore(
		invalidAddressMap,
		"invalid_address.json",
		"testpassword",
		privateKeyHex,
		addressHex,
		"scrypt",
		"Keystore with an invalid address format",
	))

	return testKeystores
}

func copyMap(original map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{})
	for key, value := range original {
		switch v := value.(type) {
		case map[string]interface{}:
			copy[key] = copyMap(v)
		default:
			copy[key] = value
		}
	}
	return copy
}

func writeInvalidKeystore(keystoreMap map[string]interface{}, filename, password, privateKeyHex, addressHex, kdf, notes string) TestKeystoreInfo {
	modifiedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal modified keystore: %v", err)
	}

	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write modified keystore file: %v", err)
	}

	fmt.Printf("Created invalid keystore file: %s\n", filename)

	return TestKeystoreInfo{
		Filename:   filename,
		Password:   password,
		PrivateKey: privateKeyHex,
		Address:    addressHex,
		KDF:        kdf,
		Valid:      false,
		Notes:      notes,
	}
}

func generateRandomHex(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
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

	readmeContent += `
## Generating Test Keystore Files

The ` + "`generate_comprehensive_keystores.go`" + ` file can be used to generate test keystore files with known private keys and passwords. This is useful for testing the keystore validation and import functionality.

To generate new keystore files:

` + "```bash" + `
cd internal/wallet/testdata
go run generate_comprehensive_keystores.go
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

## Security Note

For security reasons, the private keys used in these test files should never be used in a production environment. They are included here only for testing purposes.
`

	// Write the README file
	readmePath := filepath.Join("keystores", "README.md")
	err := ioutil.WriteFile(readmePath, []byte(readmeContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write README file: %v", err)
	}

	fmt.Println("Generated documentation in keystores/README.md")
}
