package gpg

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	commonHTTPClient "github.com/hibare/GoCommon/v2/pkg/http/client"
)

const (
	// GPGFileExtension is the extension for the GPG file.
	GPGFileExtension = "asc"

	// GPGFilePrefix is the prefix for the GPG file.
	GPGFilePrefix = "gpg_pub_key_"
)

// GPGManagerIface is the interface for the GPG manager.
type GPGManagerIface interface {
	// Read keys
	ReadPublicKeyFromFile() (string, error)
	ReadPrivateKeyFromFile() (string, error)

	// Encryption:File
	EncryptFile(inputFilePath string) (string, error)
	DecryptFile(inputFilePath string) (string, error)

	// Fetch keys
	FetchGPGPubKeyFromKeyServer(keyID, keyServerURL string) (*string, error)
}

// GPGManager is the implementation of the GPG manager.
type GPGManager struct {
	PublicKeyPath  string
	PrivateKeyPath string
	Passphrase     string
	httpClient     commonHTTPClient.HTTPClientIface
}

// ReadPublicKeyFromFile reads the public key from the file.
func (g *GPGManager) ReadPublicKeyFromFile() (string, error) {
	keyData, err := os.ReadFile(g.PublicKeyPath)
	if err != nil {
		return "", err
	}
	return string(keyData), nil
}

// ReadPrivateKeyFromFile reads the private key from the file.
func (g *GPGManager) ReadPrivateKeyFromFile() (string, error) {
	keyData, err := os.ReadFile(g.PrivateKeyPath)
	if err != nil {
		return "", err
	}
	return string(keyData), nil
}

// FetchGPGPubKeyFromKeyServer fetches a GPG key from the key server.
func (g *GPGManager) FetchGPGPubKeyFromKeyServer(keyID, keyServerURL string) (*string, error) {
	// Input validation
	if keyID == "" {
		return nil, ErrKeyIDEmpty
	}
	if keyServerURL == "" {
		return nil, ErrKeyServerURLEmpty
	}

	outputFileName := fmt.Sprintf("%s_%s.%s", GPGFilePrefix, keyID, GPGFileExtension)
	outputFilePath := filepath.Join(os.TempDir(), outputFileName)

	keyURL := fmt.Sprintf("%s/pks/lookup?op=get&search=%s", keyServerURL, keyID)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, keyURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := g.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to download GPG key: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("key-server returned non-OK status: %d", response.StatusCode)
	}

	keyData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read key data: %w", err)
	}

	file, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create key file: %w", err)
	}
	defer func() {
		if cErr := file.Close(); cErr != nil && err == nil {
			err = cErr
		}
	}()

	if _, err = file.Write(keyData); err != nil {
		return nil, fmt.Errorf("failed to write key data: %w", err)
	}

	g.PublicKeyPath = outputFilePath

	return &g.PublicKeyPath, nil
}

// EncryptFile encrypts the given file using the GPG public key and writes the result to a temp file.
func (g *GPGManager) EncryptFile(inputFilePath string) (string, error) {
	if inputFilePath == "" {
		return "", ErrEmptyInputFilePath
	}

	fileName := filepath.Base(inputFilePath)
	outputFileName := strings.TrimSuffix(fileName, fmt.Sprintf(".%s", GPGPrefix))
	outputFilePath := filepath.Join(os.TempDir(), outputFileName)

	privateKey, err := g.ReadPrivateKeyFromFile()
	if err != nil {
		return "", fmt.Errorf("failed to read private key: %w", err)
	}

	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(privateKey))
	if err != nil {
		return "", fmt.Errorf("failed to read armored key ring: %w", err)
	}
	if len(entityList) == 0 {
		return "", ErrNoEntitiesFoundInPrivateKey
	}

	entity := entityList[0]
	if entity.PrivateKey == nil {
		return "", ErrNoPrivateKeyFoundInEntity
	}

	passphraseByte := []byte(g.Passphrase)
	defer func() {
		for i := range passphraseByte {
			passphraseByte[i] = 0
		}
	}()

	if dErr := entity.PrivateKey.Decrypt(passphraseByte); dErr != nil {
		return "", fmt.Errorf("failed to decrypt private key: %w", dErr)
	}
	for _, subkey := range entity.Subkeys {
		if subkey.PrivateKey != nil {
			if dErr := subkey.PrivateKey.Decrypt(passphraseByte); dErr != nil {
				return "", fmt.Errorf("failed to decrypt subkey: %w", dErr)
			}
		}
	}

	encryptedFile, err := os.Open(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open input file: %w", err)
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

// DecryptFile decrypts the given file using the GPG private key and writes the result to a temp file.
func (g *GPGManager) DecryptFile(inputFilePath string) (string, error) {
	if inputFilePath == "" {
		return "", ErrEmptyInputFilePath
	}

	fileName := filepath.Base(inputFilePath)
	outputFileName := strings.TrimSuffix(fileName, fmt.Sprintf(".%s", GPGPrefix))
	outputFilePath := filepath.Join(os.TempDir(), outputFileName)

	privateKey, err := g.ReadPrivateKeyFromFile()
	if err != nil {
		return "", fmt.Errorf("failed to read private key: %w", err)
	}

	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(privateKey))
	if err != nil {
		return "", fmt.Errorf("failed to read armored key ring: %w", err)
	}
	if len(entityList) == 0 {
		return "", ErrNoEntitiesFoundInPrivateKey
	}

	entity := entityList[0]
	if entity.PrivateKey == nil {
		return "", ErrNoPrivateKeyFoundInEntity
	}

	passphraseByte := []byte(g.Passphrase)
	defer func() {
		for i := range passphraseByte {
			passphraseByte[i] = 0
		}
	}()

	if dErr := entity.PrivateKey.Decrypt(passphraseByte); dErr != nil {
		return "", fmt.Errorf("failed to decrypt private key: %w", dErr)
	}
	for _, subkey := range entity.Subkeys {
		if subkey.PrivateKey != nil {
			if dErr := subkey.PrivateKey.Decrypt(passphraseByte); dErr != nil {
				return "", fmt.Errorf("failed to decrypt subkey: %w", dErr)
			}
		}
	}

	encryptedFile, err := os.Open(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open input file: %w", err)
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

// ManagerOptions is the options for the GPG manager.
type ManagerOptions struct {
	HTTPClient commonHTTPClient.HTTPClientIface
}

func newGPGManager(opts ManagerOptions) GPGManagerIface {
	if opts.HTTPClient == nil {
		opts.HTTPClient = commonHTTPClient.NewDefaultClient()
	}

	return &GPGManager{
		httpClient: opts.HTTPClient,
	}
}

// NewGPGManager returns a new GPG manager.
var NewGPGManager = newGPGManager
