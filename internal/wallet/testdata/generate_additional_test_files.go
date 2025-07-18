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
)

func main() {
	// Create directory if it doesn't exist
	keystoreDir := "keystores"
	if _, err := os.Stat(keystoreDir); os.IsNotExist(err) {
		os.MkdirAll(keystoreDir, 0755)
	}

	// Generate additional test cases
	generateAdditionalTestCases()

	fmt.Println("Additional test files generated successfully!")
}

func generateAdditionalTestCases() {
	// Generate a keystore with a very long address
	generateLongAddress()

	// Generate a keystore with a malformed JSON structure
	generateMalformedJSON()

	// Generate a keystore with an empty address
	generateEmptyAddress()

	// Generate a keystore with a null address
	generateNullAddress()
}

func generateLongAddress() {
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

	// Change the address to a very long one
	keystoreMap["address"] = "5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d5d8c5d3a5e6f6d6c5b4a3a2b1c0d9e8f7a6b5c4d"

	// Write back to file
	modifiedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal modified keystore: %v", err)
	}

	filename := "long_address.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write modified keystore file: %v", err)
	}

	fmt.Printf("Created long address keystore file: %s\n", filename)

	// Update the README with the invalid keystore info
	updateInvalidKeystoreReadme(filename, "Keystore with an address that is too long")
}

func generateMalformedJSON() {
	// Create a malformed JSON file
	malformedJSON := `{
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
		},
		"extra": {
			"nested": {
				"field": "value",
				"array": [1, 2, 3, 4, 5]
			}
		}
	}`

	filename := "malformed_structure.json"
	outputPath := filepath.Join("keystores", filename)
	err := ioutil.WriteFile(outputPath, []byte(malformedJSON), 0644)
	if err != nil {
		log.Fatalf("Failed to write malformed JSON file: %v", err)
	}

	fmt.Printf("Created malformed JSON structure keystore file: %s\n", filename)

	// Update the README with the invalid keystore info
	updateInvalidKeystoreReadme(filename, "Keystore with additional unexpected fields in the JSON structure")
}

func generateEmptyAddress() {
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

	// Change the address to an empty string
	keystoreMap["address"] = ""

	// Write back to file
	modifiedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal modified keystore: %v", err)
	}

	filename := "empty_address.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write modified keystore file: %v", err)
	}

	fmt.Printf("Created empty address keystore file: %s\n", filename)

	// Update the README with the invalid keystore info
	updateInvalidKeystoreReadme(filename, "Keystore with an empty address string")
}

func generateNullAddress() {
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

	// Change the address to null
	keystoreMap["address"] = nil

	// Write back to file
	modifiedJSON, err := json.MarshalIndent(keystoreMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal modified keystore: %v", err)
	}

	filename := "null_address.json"
	outputPath := filepath.Join("keystores", filename)
	err = ioutil.WriteFile(outputPath, modifiedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write modified keystore file: %v", err)
	}

	fmt.Printf("Created null address keystore file: %s\n", filename)

	// Update the README with the invalid keystore info
	updateInvalidKeystoreReadme(filename, "Keystore with a null address value")
}

func generateRandomHex(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
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
