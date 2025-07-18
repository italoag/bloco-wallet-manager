# Additional Test Keystore Files

This file documents additional keystore files generated for testing the keystore validation and import functionality.

## Valid Keystore Files

- `real_keystore_v3_strong_scrypt.json`: Strong scrypt parameters (more secure but slower)
  - Password: `testpassword`
  - Address: `0x344995135FdECe65DFc46c673623f358E60a6626`
  - KDF: scrypt

- `real_keystore_v3_special_chars_password.json`: Password with only special characters
  - Password: `!@#$%^&*()_+{}|:<>?~`
  - Address: `0x1BBCF5946eA7197F6ccC896c85973b528812Ce3F`
  - KDF: scrypt

- `real_keystore_v3_pbkdf2_high_iterations.json`: PBKDF2 with high iteration count (20480)
  - Password: `testpassword`
  - Address: `0x4154410Ad03d0154c7b8C416ad705FE47acA7005`
  - KDF: pbkdf2


## Invalid Keystore Files

- `null_address.json`: Keystore with a null address value

- `empty_address.json`: Keystore with an empty address string

- `malformed_structure.json`: Keystore with additional unexpected fields in the JSON structure

- `long_address.json`: Keystore with an address that is too long

- `missing_kdf_params_salt.json`: Keystore missing the crypto.kdfparams.salt field

- `string_version.json`: Keystore with version as a string instead of a number

- `invalid_version_4.json`: Keystore with version 4 instead of 3

- `non_standard_cipher.json`: Keystore with non-standard cipher algorithm (aes-256-gcm)

