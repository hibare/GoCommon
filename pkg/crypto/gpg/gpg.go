// Package gpg provides utilities for GPG encryption and decryption of files.
package gpg

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
)

const (
	// GPGPrefix is the prefix for GPG files.
	GPGPrefix = "gpg"
)

type GPGIface interface {
	EncryptFile(inputFilePath string) (string, error)
	DecryptFile(inputFilePath string) (string, error)
}

// GPG holds configuration and key data for GPG operations.
type GPG struct {
	PublicKey      string
	PublicKeyPath  string
	PrivateKey     string
	PrivateKeyPath string
	Passphrase     string
}

// EncryptFile encrypts the given file using the GPG public key and writes the result to a temp file.
// Returns the path to the encrypted file on success.
func (g *GPG) EncryptFile(inputFilePath string) (string, error) {
	// Input validation
	if inputFilePath == "" {
		return "", errors.New("inputFilePath cannot be empty")
	}
	if g.PublicKey == "" {
		return "", errors.New("public key is required for encryption")
	}

	fileName := filepath.Base(inputFilePath)
	outputFileName := fmt.Sprintf("%s.%s", fileName, GPGPrefix)
	outputFilePath := filepath.Join(os.TempDir(), outputFileName)

	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(g.PublicKey))
	if err != nil {
		return "", errors.Join(errors.New("failed to read armored key ring"), err)
	}
	if len(entityList) == 0 {
		return "", errors.New("no entities found in public key")
	}

	plaintext, err := os.Open(inputFilePath)
	if err != nil {
		return "", errors.Join(errors.New("failed to open input file"), err)
	}
	defer func() {
		_ = plaintext.Close()
	}()

	output, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", errors.Join(errors.New("failed to create output file"), err)
	}
	defer func() {
		_ = output.Close()
	}()

	encrypted, err := armor.Encode(output, "PGP MESSAGE", nil)
	if err != nil {
		return "", errors.Join(errors.New("failed to create armored output"), err)
	}
	defer func() {
		_ = encrypted.Close()
	}()

	encryptionWriter, err := openpgp.Encrypt(encrypted, entityList, nil, nil, nil)
	if err != nil {
		return "", errors.Join(errors.New("failed to initialize encryption"), err)
	}
	defer func() {
		_ = encryptionWriter.Close()
	}()

	if _, err = io.Copy(encryptionWriter, plaintext); err != nil {
		return "", errors.Join(errors.New("failed to encrypt file contents"), err)
	}

	return outputFilePath, nil
}

// DecryptFile decrypts the given GPG-encrypted file using the private key and writes the result to a temp file.
// Returns the path to the decrypted file on success.
func (g *GPG) DecryptFile(inputFilePath string) (string, error) {
	// Input validation
	if inputFilePath == "" {
		return "", errors.New("inputFilePath cannot be empty")
	}
	if g.PrivateKey == "" {
		return "", errors.New("private key is required for decryption")
	}
	if g.Passphrase == "" {
		return "", errors.New("passphrase is required for decryption")
	}

	fileName := filepath.Base(inputFilePath)
	outputFileName := strings.TrimSuffix(fileName, fmt.Sprintf(".%s", GPGPrefix))
	outputFilePath := filepath.Join(os.TempDir(), outputFileName)

	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(g.PrivateKey))
	if err != nil {
		return "", errors.Join(errors.New("failed to read armored key ring"), err)
	}
	if len(entityList) == 0 {
		return "", errors.New("no entities found in private key")
	}

	entity := entityList[0]
	if entity.PrivateKey == nil {
		return "", errors.New("no private key found in entity")
	}

	passphraseByte := []byte(g.Passphrase)
	defer func() {
		for i := range passphraseByte {
			passphraseByte[i] = 0
		}
	}()

	if dErr := entity.PrivateKey.Decrypt(passphraseByte); dErr != nil {
		return "", errors.Join(errors.New("failed to decrypt private key"), dErr)
	}
	for _, subkey := range entity.Subkeys {
		if subkey.PrivateKey != nil {
			if dErr := subkey.PrivateKey.Decrypt(passphraseByte); dErr != nil {
				return "", errors.Join(errors.New("failed to decrypt subkey"), dErr)
			}
		}
	}

	encryptedFile, err := os.Open(inputFilePath)
	if err != nil {
		return "", errors.Join(errors.New("failed to open input file"), err)
	}
	defer func() {
		_ = encryptedFile.Close()
	}()

	decoded, err := armor.Decode(encryptedFile)
	if err != nil {
		return "", errors.Join(errors.New("failed to decode armored input"), err)
	}

	md, err := openpgp.ReadMessage(decoded.Body, entityList, nil, nil)
	if err != nil {
		return "", errors.Join(errors.New("failed to read PGP message"), err)
	}

	outputFile, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", errors.Join(errors.New("failed to create output file"), err)
	}
	defer func() {
		_ = outputFile.Close()
	}()

	if _, err = io.Copy(outputFile, md.UnverifiedBody); err != nil {
		return "", errors.Join(errors.New("failed to write decrypted contents"), err)
	}

	return outputFilePath, nil
}

// Options is the options for the GPG service.
type Options struct {
	PublicKey      string
	PublicKeyPath  string
	PrivateKey     string
	PrivateKeyPath string
	Passphrase     string
}

func newGPG(opts Options) GPGIface {
	return &GPG{
		PublicKey:      opts.PublicKey,
		PublicKeyPath:  opts.PublicKeyPath,
		PrivateKey:     opts.PrivateKey,
		PrivateKeyPath: opts.PrivateKeyPath,
		Passphrase:     opts.Passphrase,
	}
}

// NewGPG returns a new GPG instance.
var NewGPG = newGPG
