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

func main() {
	// Create directory if it doesn't exist
	keystoreDir := "keystores"
	if _, err := os.Stat(keystoreDir); os.IsNotExist(err) {
		os.MkdirAll(keystoreDir, 0755)
	}

	// Generate valid keystores
	generateValidKeystores()

	// Generate invalid keystores
	generateInvalidKeystores()

	fmt.Println("All test keystore files generated successfully!")
}

func generateValidKeystores() {
	// Generate keystore with strong scrypt parameters
	generateRealKeystore(
		"real_keystore_v3_strong_scrypt.json",
		"testpassword",
		keystore.StandardScryptN*2, // Double the standard parameters
		keystore.StandardScryptP*2,
		"Strong scrypt parameters (more secure but slower)",
	)

	// Generate keystore with special characters in the password
	generateRealKeystore(
		"real_keystore_v3_special_chars_password.json",
		"!@#$%^&*()_+{}|:<>?~",
		keystore.StandardScryptN,
		keystore.StandardScryptP,
		"Password with only special characters",
	)

	// Generate keystore with PBKDF2 and high iteration count
	generatePBKDF2Keystore(
		"real_keystore_v3_pbkdf2_high_iterations.json",
		"testpassword",
		20480, // Double the standard iterations
		"PBKDF2 with high iteration count (20480)",
	)
}

func generateInvalidKeystores() {
	// Generate keystore with non-standard cipher
	generateNonStandardCipher()

	// Generate keystore with invalid version
	generateInvalidVersion()

	// Generate keystore with string version
	generateStringVersion()

	// Generate keystore with missing fields
	generateMissingFields()
}

func generateRealKeystore(filename, password string, scryptN, scryptP int, notes string) {
	// Generate a new random private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
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

	// Update the README with the password and address
	updateReadme(filename, password, addressHex, "scrypt", notes)
}

func generatePBKDF2Keystore(filename, password string, iterations int, notes string) {
	// Generate a new random private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
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

	// Update the README with the password and address
	updateReadme(filename, password, addressHex, "pbkdf2", notes)
}

func generateNonStandardCipher() {
	// Generate a new random private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	_ = address.Hex() // Not used but included for consistency

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

	// Update the README with the invalid keystore info
	updateInvalidKeystoreReadme(filename, "Keystore with non-standard cipher algorithm (aes-256-gcm)")
}

func generateInvalidVersion() {
	// Generate a new random private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Get the address (not used but included for consistency)
	_ = crypto.PubkeyToAddress(privateKey.PublicKey)

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
	keystoreMap["version"] = 4 // Only version 3 is valid

	// Write back to file
	modifiedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal modified keystore: %v", err)
	}

	filename := "invalid_version_4.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write modified keystore file: %v", err)
	}

	fmt.Printf("Created invalid version keystore file: %s\n", filename)

	// Update the README with the invalid keystore info
	updateInvalidKeystoreReadme(filename, "Keystore with version 4 instead of 3")
}

func generateStringVersion() {
	// Generate a new random private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Get the address (not used but included for consistency)
	_ = crypto.PubkeyToAddress(privateKey.PublicKey)

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

	// Update the README with the invalid keystore info
	updateInvalidKeystoreReadme(filename, "Keystore with version as a string instead of a number")
}

func generateMissingFields() {
	// Generate a new random private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

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

	// Create a copy for each missing field

	// Missing crypto.kdfparams.salt
	missingKDFParamsSaltMap := copyMap(keystoreMap)
	crypto := missingKDFParamsSaltMap["crypto"].(map[string]interface{})
	kdfparams := crypto["kdfparams"].(map[string]interface{})
	delete(kdfparams, "salt")

	filename := "missing_kdf_params_salt.json"
	outputPath := filepath.Join("keystores", filename)
	modifiedJSON, _ := json.MarshalIndent(missingKDFParamsSaltMap, "", "  ")
	ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	fmt.Printf("Created missing KDF params salt keystore file: %s\n", filename)
	updateInvalidKeystoreReadme(filename, "Keystore missing the crypto.kdfparams.salt field")
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

func generateRandomHex(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func updateReadme(filename, password, address, kdf, notes string) {
	// Create or append to the ADDITIONAL_KEYSTORES.md file
	readmePath := filepath.Join("keystores", "ADDITIONAL_KEYSTORES.md")

	var content string
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		// Create new file with header
		content = "# Additional Test Keystore Files\n\n"
		content += "This file documents additional keystore files generated for testing the keystore validation and import functionality.\n\n"
		content += "## Valid Keystore Files\n\n"
	} else {
		// Read existing content
		existingContent, err := ioutil.ReadFile(readmePath)
		if err != nil {
			log.Fatalf("Failed to read README: %v", err)
		}
		content = string(existingContent)

		// Check if we need to add the Valid Keystore Files section
		if !strings.Contains(content, "## Valid Keystore Files") {
			content += "\n## Valid Keystore Files\n\n"
		}
	}

	// Add the keystore information
	entry := fmt.Sprintf("- `%s`: %s\n", filename, notes)
	entry += fmt.Sprintf("  - Password: `%s`\n", password)
	entry += fmt.Sprintf("  - Address: `%s`\n", address)
	entry += fmt.Sprintf("  - KDF: %s\n\n", kdf)

	// Check if this entry already exists
	if !strings.Contains(content, fmt.Sprintf("- `%s`:", filename)) {
		// Find the position to insert the new entry
		validSection := strings.Index(content, "## Valid Keystore Files")
		invalidSection := strings.Index(content, "## Invalid Keystore Files")

		if invalidSection > validSection && invalidSection != -1 {
			// Insert before the Invalid section
			content = content[:invalidSection] + entry + content[invalidSection:]
		} else {
			// Append to the end
			content += entry
		}

		// Write the updated content
		err := ioutil.WriteFile(readmePath, []byte(content), 0644)
		if err != nil {
			log.Fatalf("Failed to update README: %v", err)
		}
	}
}

func updateInvalidKeystoreReadme(filename, notes string) {
	// Create or append to the ADDITIONAL_KEYSTORES.md file
	readmePath := filepath.Join("keystores", "ADDITIONAL_KEYSTORES.md")

	var content string
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		// Create new file with header
		content = "# Additional Test Keystore Files\n\n"
		content += "This file documents additional keystore files generated for testing the keystore validation and import functionality.\n\n"
		content += "## Invalid Keystore Files\n\n"
	} else {
		// Read existing content
		existingContent, err := ioutil.ReadFile(readmePath)
		if err != nil {
			log.Fatalf("Failed to read README: %v", err)
		}
		content = string(existingContent)

		// Check if we need to add the Invalid Keystore Files section
		if !strings.Contains(content, "## Invalid Keystore Files") {
			content += "\n## Invalid Keystore Files\n\n"
		}
	}

	// Add the keystore information
	entry := fmt.Sprintf("- `%s`: %s\n\n", filename, notes)

	// Check if this entry already exists
	if !strings.Contains(content, fmt.Sprintf("- `%s`:", filename)) {
		// Find the position to insert the new entry
		invalidSection := strings.Index(content, "## Invalid Keystore Files")

		if invalidSection != -1 {
			// Find the end of the section header line
			endOfHeader := invalidSection + len("## Invalid Keystore Files\n\n")
			content = content[:endOfHeader] + entry + content[endOfHeader:]
		} else {
			// Add the section and entry
			content += "\n## Invalid Keystore Files\n\n" + entry
		}

		// Write the updated content
		err := ioutil.WriteFile(readmePath, []byte(content), 0644)
		if err != nil {
			log.Fatalf("Failed to update README: %v", err)
		}
	}
}
