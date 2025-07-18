# BLOCO Wallet Manager

BLOCO Wallet Manager is a command-line interface (CLI) application designed to manage cryptocurrency wallets compatible with the Ethereum network and adhering to the KeyStoreV3 standard.

## Core Features

- **Wallet Management**: Create, import, export, delete, block, and unblock wallet addresses
- **Security**: KeyStoreV3 format compatibility with plans for external vault integration
- **Balance Inquiry**: Query balances of Ethereum-compatible wallets
- **Terminal User Interface (TUI)**: Intuitive interface for wallet management

## Current Development Focus

The project is currently focused on improving keystore validation and import functionality, with specific attention to:

1. Enhanced validation of KeyStoreV3 format files
2. Proper error handling with localized error messages
3. Deterministic mnemonic generation from private keys
4. Comprehensive testing of wallet import flows

## Future Roadmap

- Integration with external vaults (Hashicorp Vault, AWS KMS, Cloud HSM, Azure Key Vault)
- Support for additional blockchain networks
- Advanced cryptography support (secp256r1, ed25519 curves; ECDSA, EdDSA algorithms)
- Enhanced security features (2FA, multi-signature wallets)