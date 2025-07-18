# Keystore Files

This directory contains keystore files for the BLOCO Wallet Manager. These files follow the KeystoreV3 standard used by Ethereum wallets.

## Flexible Keystore Import

The BLOCO Wallet Manager now supports importing keystore files with any file extension, not just `.json` files. This means you can import keystores with extensions like:
- `.key`
- `.keystoremaster`
- `.eth`
- No extension at all

## File Format

Regardless of the file extension, the content must be a valid KeystoreV3 JSON structure. The system validates the content of the file rather than relying on the file extension.

If the file content is invalid, you'll receive one of these error messages:
- "O arquivo não contém um JSON válido" (The file does not contain valid JSON)
- "O arquivo não contém um keystore v3 válido" (The file does not contain a valid keystore v3)

## Example Files

This directory contains several example keystore files:

- `real_keystore_v3_complex_password.json`: A keystore with a complex password
- `real_keystore_v3_empty_password.json`: A keystore with an empty password
- `real_keystore_v3_light.json`: A keystore with light scrypt parameters
- `real_keystore_v3_pbkdf2_proper.json`: A keystore using PBKDF2 KDF
- `real_keystore_v3_standard.json`: A keystore with standard scrypt parameters

## Security Note

For security reasons, never share your keystore files or passwords with anyone. The keystores in this directory are for testing purposes only and should not be used in production environments.