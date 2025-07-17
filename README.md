<!--suppress HtmlDeprecatedAttribute -->

<h1 align="center">🔐 BLOCO | Wallet Manager</h1>

<p align="center">
<a href="https://github.com/italoag/bloco-wallet-manager/releases" style="text-decoration: none">
<img src="https://img.shields.io/github/v/release/italoag/bloco-wallet-manager?style=flat-square" alt="Latest Release">
</a>

<a href="https://github.com/italoag/bloco-wallet-manager/stargazers" style="text-decoration: none">
<img src="https://img.shields.io/github/stars/italoag/bloco-wallet-manager.svg?style=flat-square" alt="Stars">
</a>

<a href="https://github.com/italoag/bloco-wallet-manager/fork" style="text-decoration: none">
<img src="https://img.shields.io/github/forks/italoag/bloco-wallet-manager.svg?style=flat-square" alt="Forks">
</a>

<a href="https://opensource.org/licenses/MIT" style="text-decoration: none">
<img src="https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square" alt="License: MIT">
</a>



<br/>

<a href="https://github.com/italoag/bloco-wallet-manager/releases" style="text-decoration: none">
<img src="https://img.shields.io/badge/platform-windows%20%7C%20macos%20%7C%20linux-informational?style=for-the-badge" alt="Downloads">
</a>

 <a href="https://twitter.com/0xitalo">
        <img src="https://img.shields.io/badge/Twitter-%400xitalo-1DA1F2?logo=twitter&style=for-the-badge" alt=""/>
    </a>

<br/>
<br/>

### Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

### Introduction
**BLOCO Wallet Manager** is a command-line interface (CLI) application designed to manage cryptocurrency wallets compatible with the Ethereum network and adhering to the KeyStoreV3 standard. Developed in GoLang, BLOCO provides a Terminal User Interface (TUI) for seamless wallet management. Future integrations will include external vaults such as Hashicorp Vault, Amazon KMS, Cloud HSM, and Azure Key Vault.

### Features
- **Terminal User Interface**
  - Navigable TUI built with Bubble Tea.
  - Multi-language support (English, Portuguese and Spanish).
  - Splash screen rendered with random ASCII fonts.
- **Wallet Operations**
  - Create new Ethereum wallets secured by password.
  - Import wallets from mnemonic phrases or raw private keys.
  - View wallet details after password verification.
  - List and delete stored wallets.
- **Persistence & Configuration**
  - Wallet metadata stored in a local SQLite database.
  - Keystore files saved in a configurable directory using KeyStoreV3.
  - Application settings and fonts managed via YAML and JSON files.
  - Logging to `blocowallet.log` for troubleshooting.

### Installation
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

### Usage
Run the BLOCO Wallet Manager using the terminal:

```bash
blocowm
```

Navigate through the TUI to manage your wallets. Available commands include:

- **Create Wallet:** Generate a new Ethereum wallet.
- **Import from Mnemonic:** Restore a wallet using a mnemonic phrase.
- **Import from Private Key:** Load a wallet from a raw private key.
- **List Wallets:** Display stored wallets and view details or delete them.
### Roadmap
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

### Contributing
Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes with clear messages.
4. Submit a pull request detailing your changes.


### License
This project is licensed under the [MIT License](LICENSE).
