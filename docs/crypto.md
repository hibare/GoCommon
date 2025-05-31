# Crypto Package Documentation

## Overview

The `crypto` package provides utilities for cryptographic operations. The main subpackage is `gpg`, which offers GPG encryption and decryption for files, as well as key management utilities.

---

## Subpackages

- **gpg**: Utilities for GPG encryption/decryption, key download, and file operations.

---

## Key Types and Functions (GPG)

- **GPG**: Struct holding configuration and key data for GPG operations.
- **DownloadGPGPubKey(keyID, keyServerURL)**: Downloads a GPG public key from a key server.
- **EncryptFile(inputFilePath)**: Encrypts a file using the GPG public key.
- **DecryptFile(inputFilePath)**: Decrypts a GPG-encrypted file using the private key and passphrase.

---

## Example Usage

```go
import (
    "github.com/hibare/GoCommon/v2/pkg/crypto/gpg"
)

gpgKey, err := gpg.DownloadGPGPubKey("keyid", "https://keyserver.example.com")
if err != nil {
    panic(err)
}
gpgKey.Passphrase = "your-passphrase"
encFile, err := gpgKey.EncryptFile("/path/to/file.txt")
// ...
```

---

## Notes

- Uses [ProtonMail/go-crypto](https://github.com/ProtonMail/go-crypto) for OpenPGP operations.
- Designed for secure file encryption/decryption and key management.
