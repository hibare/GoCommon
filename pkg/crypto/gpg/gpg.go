// Package gpg provides utilities for GPG encryption and decryption of files.
package gpg

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/hibare/GoCommon/v2/pkg/errors"
)

// GPG holds configuration and key data for GPG operations.
type GPG struct {
	KeyID          string
	KeyServerURL   string
	PublicKey      string
	PublicKeyPath  string
	PrivateKey     string
	PrivateKeyPath string
	Passphrase     string
}

const GPGPrefix = "gpg"

// DownloadGPGPubKey downloads a GPG public key from a key server and saves it to a temp file.
func DownloadGPGPubKey(keyID, keyServerURL string) (GPG, error) {
	gpgPubKey := GPG{
		KeyID:        keyID,
		KeyServerURL: keyServerURL,
	}

	keyURL := fmt.Sprintf("%s/pks/lookup?op=get&search=%s", keyServerURL, keyID)
	response, err := http.Get(keyURL)
	if err != nil {
		return gpgPubKey, fmt.Errorf("failed to download GPG key: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return gpgPubKey, fmt.Errorf("key-server returned non-OK status: %w", errors.ErrNonOKError)
	}

	keyData, err := io.ReadAll(response.Body)
	if err != nil {
		return gpgPubKey, fmt.Errorf("failed to read key data: %w", err)
	}

	outputFileName := fmt.Sprintf("gpg_pub_key_%s.asc", keyID)
	outputFilePath := filepath.Join(os.TempDir(), outputFileName)

	file, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return gpgPubKey, fmt.Errorf("failed to create key file: %w", err)
	}
	defer func() {
		if cErr := file.Close(); cErr != nil && err == nil {
			err = cErr
		}
	}()

	if _, err = file.Write(keyData); err != nil {
		return gpgPubKey, fmt.Errorf("failed to write key data: %w", err)
	}

	gpgPubKey.PublicKeyPath = outputFilePath
	gpgPubKey.PublicKey = string(keyData)

	return gpgPubKey, nil
}

// EncryptFile encrypts the given file using the GPG public key and writes the result to a temp file.
// Returns the path to the encrypted file on success.
func (g *GPG) EncryptFile(inputFilePath string) (string, error) {
	fileName := filepath.Base(inputFilePath)
	outputFileName := fmt.Sprintf("%s.%s", fileName, GPGPrefix)
	outputFilePath := filepath.Join(os.TempDir(), outputFileName)

	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(g.PublicKey))
	if err != nil {
		return "", fmt.Errorf("failed to read armored key ring: %w", err)
	}
	if len(entityList) == 0 {
		return "", fmt.Errorf("no entities found in public key")
	}

	plaintext, err := os.Open(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open input file %q: %w", inputFilePath, err)
	}
	defer func() {
		_ = plaintext.Close()
	}()

	output, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		_ = output.Close()
	}()

	encrypted, err := armor.Encode(output, "PGP MESSAGE", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create armored output: %w", err)
	}
	defer func() {
		_ = encrypted.Close()
	}()

	encryptionWriter, err := openpgp.Encrypt(encrypted, entityList, nil, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to initialize encryption: %w", err)
	}
	defer func() {
		_ = encryptionWriter.Close()
	}()

	if _, err = io.Copy(encryptionWriter, plaintext); err != nil {
		return "", fmt.Errorf("failed to encrypt file contents: %w", err)
	}

	return outputFilePath, nil
}

// DecryptFile decrypts the given GPG-encrypted file using the private key and writes the result to a temp file.
// Returns the path to the decrypted file on success.
func (g *GPG) DecryptFile(inputFilePath string) (string, error) {
	fileName := filepath.Base(inputFilePath)
	outputFileName := strings.TrimSuffix(fileName, fmt.Sprintf(".%s", GPGPrefix))
	outputFilePath := filepath.Join(os.TempDir(), outputFileName)

	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(g.PrivateKey))
	if err != nil {
		return "", fmt.Errorf("failed to read armored key ring: %w", err)
	}
	if len(entityList) == 0 {
		return "", fmt.Errorf("no entities found in private key")
	}

	entity := entityList[0]
	if entity.PrivateKey == nil {
		return "", fmt.Errorf("no private key found in entity")
	}

	passphraseByte := []byte(g.Passphrase)
	defer func() {
		for i := range passphraseByte {
			passphraseByte[i] = 0
		}
	}()

	if err := entity.PrivateKey.Decrypt(passphraseByte); err != nil {
		return "", fmt.Errorf("failed to decrypt private key: %w", err)
	}
	for _, subkey := range entity.Subkeys {
		if subkey.PrivateKey != nil {
			if err := subkey.PrivateKey.Decrypt(passphraseByte); err != nil {
				return "", fmt.Errorf("failed to decrypt subkey: %w", err)
			}
		}
	}

	encryptedFile, err := os.Open(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open input file %q: %w", inputFilePath, err)
	}
	defer func() {
		_ = encryptedFile.Close()
	}()

	decoded, err := armor.Decode(encryptedFile)
	if err != nil {
		return "", fmt.Errorf("failed to decode armored input: %w", err)
	}

	md, err := openpgp.ReadMessage(decoded.Body, entityList, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to read PGP message: %w", err)
	}

	outputFile, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		_ = outputFile.Close()
	}()

	if _, err = io.Copy(outputFile, md.UnverifiedBody); err != nil {
		return "", fmt.Errorf("failed to write decrypted contents: %w", err)
	}

	return outputFilePath, nil
}
