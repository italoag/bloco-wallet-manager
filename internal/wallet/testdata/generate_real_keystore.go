package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Create a known private key (32 bytes = 64 hex characters)
	privateKeyHex := "1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b"
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
	fmt.Printf("Address: %s\n", address.Hex())

	// Create a temporary directory for the keystore
	tempDir, err := ioutil.TempDir("", "keystore")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a keystore with scrypt
	scryptKs := keystore.NewKeyStore(tempDir, keystore.StandardScryptN, keystore.StandardScryptP)
	password := "testpassword"

	// Import the private key
	account, err := scryptKs.ImportECDSA(privateKey, password)
	if err != nil {
		log.Fatalf("Failed to import private key: %v", err)
	}

	// Read the keystore file
	keystoreJSON, err := ioutil.ReadFile(account.URL.Path)
	if err != nil {
		log.Fatalf("Failed to read keystore file: %v", err)
	}

	// Write the keystore file to the testdata directory
	err = ioutil.WriteFile("keystores/real_keystore_v3.json", keystoreJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write keystore file: %v", err)
	}

	fmt.Println("Created real keystore file with:")
	fmt.Printf("Private Key: %s\n", privateKeyHex)
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Address: %s\n", address.Hex())

	// Create a keystore with pbkdf2
	pbkdf2Dir, err := ioutil.TempDir("", "keystore-pbkdf2")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(pbkdf2Dir)

	pbkdf2Ks := keystore.NewKeyStore(pbkdf2Dir, keystore.LightScryptN, keystore.LightScryptP)

	// Import the private key
	pbkdf2Account, err := pbkdf2Ks.ImportECDSA(privateKey, password)
	if err != nil {
		log.Fatalf("Failed to import private key: %v", err)
	}

	// Read the keystore file
	pbkdf2KeystoreJSON, err := ioutil.ReadFile(pbkdf2Account.URL.Path)
	if err != nil {
		log.Fatalf("Failed to read keystore file: %v", err)
	}

	// Write the keystore file to the testdata directory
	err = ioutil.WriteFile("keystores/real_keystore_v3_pbkdf2.json", pbkdf2KeystoreJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write keystore file: %v", err)
	}

	fmt.Println("Created real PBKDF2 keystore file with the same credentials")
}
