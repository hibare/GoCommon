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
	"github.com/hibare/GoCommon/pkg/errors"
)

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

func DownloadGPGPubKey(keyID, keyServerURL string) (GPG, error) {
	gpgPubKey := GPG{
		KeyID:        keyID,
		KeyServerURL: keyServerURL,
	}

	keyURL := fmt.Sprintf("%s/pks/lookup?op=get&search=%s", keyServerURL, keyID)
	response, err := http.Get(keyURL)
	if err != nil {
		return gpgPubKey, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return gpgPubKey, errors.ErrNonOKError
	}

	// Save the GPG key to a file
	keyData, err := io.ReadAll(response.Body)
	if err != nil {
		return gpgPubKey, err
	}

	// Create a file in temp dir
	outputFileName := fmt.Sprintf("gpg_pub_key_%s.asc", keyID)
	outputFilePath := filepath.Join(os.TempDir(), outputFileName)

	// Create or open the file for writing (with more control over file permissions)
	file, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return gpgPubKey, err
	}
	defer file.Close()

	// Write data to the file
	_, err = file.Write(keyData)
	if err != nil {
		return gpgPubKey, err
	}

	gpgPubKey.PublicKeyPath = outputFilePath
	gpgPubKey.PublicKey = string(keyData)

	return gpgPubKey, nil
}

func (g *GPG) EncryptFile(inputFilePath string) (string, error) {
	fileName := filepath.Base(inputFilePath)

	// Create output file in temp dir
	outputFileName := fmt.Sprintf("%s.%s", fileName, GPGPrefix)
	outputFilePath := filepath.Join(os.TempDir(), outputFileName)

	// Create an entity list from the public key
	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(g.PublicKey))
	if err != nil {
		return outputFilePath, fmt.Errorf("failed to read armored key ring: %s", err)
	}

	// Encrypt the file using the public key
	plaintext, err := os.Open(inputFilePath)
	if err != nil {
		return outputFilePath, fmt.Errorf("failed to open input file: %s", err)
	}
	defer plaintext.Close()

	// Create the output file
	output, err := os.Create(outputFilePath)
	if err != nil {
		return outputFilePath, fmt.Errorf("failed to create output file: %s", err)
	}
	defer output.Close()

	encrypted, err := armor.Encode(output, "PGP MESSAGE", nil)
	if err != nil {
		return outputFilePath, fmt.Errorf("failed to create armored output: %s", err)
	}
	defer encrypted.Close()

	encryptionWriter, err := openpgp.Encrypt(encrypted, entityList, nil, nil, nil)
	if err != nil {
		return outputFilePath, fmt.Errorf("failed to initialize encryption: %s", err)
	}

	// Copy the contents of the plaintext file to the encryption writer
	_, err = io.Copy(encryptionWriter, plaintext)
	if err != nil {
		return outputFilePath, fmt.Errorf("failed to encrypt file contents: %s", err)
	}

	encryptionWriter.Close()

	return outputFilePath, nil
}

func (g *GPG) DecryptFile(inputFilePath string) (string, error) {
	fileName := filepath.Base(inputFilePath)

	// Create output file in temp dir
	outputFileName := strings.TrimSuffix(fileName, fmt.Sprintf(".%s", GPGPrefix))
	outputFilePath := filepath.Join(os.TempDir(), outputFileName)

	// Create an entity list from the private key
	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(g.PrivateKey))
	if err != nil {
		return outputFilePath, fmt.Errorf("failed to read armored key ring: %s", err)
	}

	entity := entityList[0]

	// Get the passphrase and read the private key.
	passphraseByte := []byte(g.Passphrase)
	entity.PrivateKey.Decrypt(passphraseByte)
	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphraseByte)
	}

	// Read the encrypted file
	encryptedFile, err := os.Open(inputFilePath)
	if err != nil {
		return outputFilePath, fmt.Errorf("failed to open input file: %s", err)
	}
	defer encryptedFile.Close()

	decoded, err := armor.Decode(encryptedFile)
	if err != nil {
		return outputFilePath, fmt.Errorf("failed to decode armored input: %s", err)
	}

	md, err := openpgp.ReadMessage(decoded.Body, entityList, nil, nil)
	if err != nil {
		return outputFilePath, err
	}

	// Create the output file
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return outputFilePath, fmt.Errorf("failed to create output file: %s", err)
	}
	defer outputFile.Close()

	// Decrypt and write the contents to the output file
	_, err = outputFile.ReadFrom(md.UnverifiedBody)
	return outputFilePath, err
}
