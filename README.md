# BLOCO Wallet Manager

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

## Introduction
**BLOCO Wallet Manager** is a command-line interface (CLI) application designed to manage cryptocurrency wallets compatible with the Ethereum network and adhering to the KeyStoreV3 standard. Developed in GoLang, BLOCO provides a Terminal User Interface (TUI) for seamless wallet management. Future integrations will include external vaults such as Hashicorp Vault, Amazon KMS, Cloud HSM, and Azure Key Vault.

## Features
- **Wallet Management**
    - Create new wallets compatible with Ethereum.
    - Import wallets using Mnemonics.
    - Export wallets in KeyStoreV3 format.
    - Delete, block, and unblock wallet addresses.
    - List all managed wallets.

- **Security**
    - Compatibility with KeyStoreV3 for secure key storage.
    - Planned integration with external vaults:
        - Hashicorp Vault
        - Amazon KMS
        - Cloud HSM
        - Azure Key Vault

- **Balance Inquiry**
    - Query the balance of Ethereum-compatible wallets.

- **Extensibility**
    - Support for additional blockchain networks (future).
    - Support for multiple cryptographic curves and signature algorithms:
        - Curves: secp256k1, secp256r1, ed25519
        - Algorithms: ECDSA, EdDSA

## Installation
Ensure you have [Go](https://golang.org/doc/install) installed on your system.

```bash
git clone https://github.com/italoag/go-wallet-manager.git
cd go-wallet-manager
go build -o blocowm
```

Move the executable to a directory in your PATH for easy access:

```bash
mv blocowm /usr/local/bin/
```

## Usage
Run the BLOCO Wallet Manager using the terminal:

```bash
blocowm
```

Navigate through the TUI to manage your wallets. Available commands include:

- **Create Wallet:** Initialize a new Ethereum-compatible wallet.
- **Import Wallet:** Import existing wallets using Mnemonics.
- **List Wallets:** Display all managed wallets.

## Roadmap
**Upcoming Features:**

- **Vault Integrations:**
    - Hashicorp Vault
    - Amazon KMS
    - Cloud HSM
    - Azure Key Vault

- **Multi-Network Support:**
    - Integration with additional blockchain networks.

- **Advanced Cryptography:**
    - Support for secp256r1 and ed25519 curves.
    - Implementation of ECDSA and EdDSA signature algorithms.

- **Enhanced Security Features:**
    - **Import Wallet:** Import existing wallets using private keys.
    - Two-factor authentication for wallet access.
    - Multi-signature wallet support.

- **User Experience Improvements:**
    - Enhanced TUI with more intuitive navigation.
    - Detailed transaction histories and analytics.
    - Delete an account address.
    - **Delete Wallet:** Remove a wallet from the manager.
    - **Block/Unblock Wallet:** Temporarily disable or enable a wallet address.
    - **Check Balance:** View the balance of a selected wallet.

## Contributing
Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes with clear messages.
4. Submit a pull request detailing your changes.

For more information, refer to the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## License
This project is licensed under the [MIT License](LICENSE).
