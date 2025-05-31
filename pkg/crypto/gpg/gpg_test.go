package gpg

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testPublicKey = "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\n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\nKtQIWXTIGw==\n=oPyT\n-----END PGP PUBLIC KEY BLOCK-----"
	//nolint: gosec  // Ignoring false positive "G101 Potential hardcoded credentials"
	testPrivateKey  = "-----BEGIN PGP PRIVATE KEY BLOCK-----\n\n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\n=ZxHp\n-----END PGP PRIVATE KEY BLOCK-----"
	testPassphrase  = "password"
	testFileContent = "This is a sample text file created in the temp directory."
)

func setupTestFile(t *testing.T) string {
	t.Helper()
	tempFile, err := os.CreateTemp(t.TempDir(), "sample*.txt")
	require.NoError(t, err)
	_, err = tempFile.WriteString(testFileContent)
	require.NoError(t, err)
	require.NoError(t, tempFile.Close())
	return tempFile.Name()
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestDownloadGPGPubKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, testPublicKey)
	}))
	t.Cleanup(server.Close)

	keyServerURL := server.URL
	keyID := "dummy_key_id"

	gpgPubKey, err := DownloadGPGPubKey(keyID, keyServerURL)
	require.NoError(t, err)

	_, err = os.Stat(gpgPubKey.PublicKeyPath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Remove(gpgPubKey.PublicKeyPath) })

	downloadedData, err := os.ReadFile(gpgPubKey.PublicKeyPath)
	require.NoError(t, err)
	require.NotEmpty(t, string(downloadedData))
	require.Equal(t, testPublicKey, string(downloadedData))
	require.Equal(t, testPublicKey, gpgPubKey.PublicKey)
}

func TestDownloadGPGPubKeyNoKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	keyServerURL := server.URL
	keyID := "dummy_key_id"

	gpgPubKey, err := DownloadGPGPubKey(keyID, keyServerURL)
	require.Error(t, err)
	require.Contains(t, err.Error(), "key-server returned non-OK status")
	require.Empty(t, gpgPubKey.PublicKey)
	require.Empty(t, gpgPubKey.PublicKeyPath)
}

func TestGPGEncrypt(t *testing.T) {
	tempFilePath := setupTestFile(t)
	gpgPubkey := GPG{PublicKey: testPublicKey}

	encryptedFilePath, err := gpgPubkey.EncryptFile(tempFilePath)
	require.NoError(t, err)
	_, err = os.Stat(encryptedFilePath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Remove(encryptedFilePath) })
}

func TestGPGEncryptInvalidPubKey(t *testing.T) {
	tempFilePath := setupTestFile(t)
	gpgPubkey := GPG{PublicKey: "invalid key"}

	encryptedFilePath, err := gpgPubkey.EncryptFile(tempFilePath)
	require.Error(t, err)
	require.Empty(t, encryptedFilePath)
}

func TestGPGEncryptInvalidFile(t *testing.T) {
	tempFilePath := "/tmp/non-exists-file.txt"
	gpgPubkey := GPG{PublicKey: testPublicKey}

	encryptedFilePath, err := gpgPubkey.EncryptFile(tempFilePath)
	require.Error(t, err)
	require.Empty(t, encryptedFilePath)
}

func TestGPGDecrypt(t *testing.T) {
	tempFilePath := setupTestFile(t)
	gpgPubkey := GPG{
		PublicKey:  testPublicKey,
		PrivateKey: testPrivateKey,
		Passphrase: testPassphrase,
	}

	encryptedFilePath, err := gpgPubkey.EncryptFile(tempFilePath)
	require.NoError(t, err)
	_, err = os.Stat(encryptedFilePath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Remove(encryptedFilePath) })

	decryptedFilePath, err := gpgPubkey.DecryptFile(encryptedFilePath)
	require.NoError(t, err)
	_, err = os.Stat(decryptedFilePath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Remove(decryptedFilePath) })

	decryptedData, err := os.ReadFile(decryptedFilePath)
	require.NoError(t, err)
	require.Equal(t, testFileContent, string(decryptedData))
}

func TestGPGDecryptInvalidPass(t *testing.T) {
	tempFilePath := setupTestFile(t)
	gpgPubkey := GPG{
		PublicKey:  testPublicKey,
		PrivateKey: testPrivateKey,
		Passphrase: "invalid",
	}

	encryptedFilePath, err := gpgPubkey.EncryptFile(tempFilePath)
	require.NoError(t, err)
	_, err = os.Stat(encryptedFilePath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Remove(encryptedFilePath) })

	decryptedFilePath, err := gpgPubkey.DecryptFile(encryptedFilePath)
	require.Error(t, err)
	require.Empty(t, decryptedFilePath)
}

func TestGPGDecryptInvalidFile(t *testing.T) {
	encryptedFilePath := "/tmp/non-exists-file.txt.gpg"
	gpgPubkey := GPG{
		PublicKey:  testPublicKey,
		PrivateKey: testPrivateKey,
		Passphrase: "invalid",
	}
	decryptedFilePath, err := gpgPubkey.DecryptFile(encryptedFilePath)
	require.Error(t, err)
	require.Empty(t, decryptedFilePath)
}

func TestGPGDecryptInvalidPrivKey(t *testing.T) {
	encryptedFilePath := "/tmp/non-exists-file.txt.gpg"
	gpgPubkey := GPG{
		PrivateKey: "invalid",
		Passphrase: "invalid",
	}
	decryptedFilePath, err := gpgPubkey.DecryptFile(encryptedFilePath)
	require.Error(t, err)
	require.Empty(t, decryptedFilePath)
}

func TestGPGEncryptEmptyPublicKey(t *testing.T) {
	tempFilePath := setupTestFile(t)
	gpgPubkey := GPG{PublicKey: ""}

	encryptedFilePath, err := gpgPubkey.EncryptFile(tempFilePath)
	require.Error(t, err)
	require.Empty(t, encryptedFilePath)
}

func TestGPGDecryptEmptyPrivateKey(t *testing.T) {
	tempFilePath := setupTestFile(t)
	gpgPubkey := GPG{
		PublicKey:  testPublicKey,
		PrivateKey: "",
		Passphrase: testPassphrase,
	}

	encryptedFilePath, err := gpgPubkey.EncryptFile(tempFilePath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Remove(encryptedFilePath) })

	decryptedFilePath, err := gpgPubkey.DecryptFile(encryptedFilePath)
	require.Error(t, err)
	require.Empty(t, decryptedFilePath)
}

func TestGPGEncryptEmptyInputFile(t *testing.T) {
	tempFile, err := os.CreateTemp(t.TempDir(), "emptyfile*.txt")
	require.NoError(t, err)
	tempFilePath := tempFile.Name()
	require.NoError(t, tempFile.Close())

	gpgPubkey := GPG{PublicKey: testPublicKey}

	encryptedFilePath, err := gpgPubkey.EncryptFile(tempFilePath)
	require.NoError(t, err)
	_, err = os.Stat(encryptedFilePath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Remove(encryptedFilePath) })
}

func TestGPGDecryptNotPGPFile(t *testing.T) {
	tempFilePath := setupTestFile(t)
	gpgPubkey := GPG{
		PrivateKey: testPrivateKey,
		Passphrase: testPassphrase,
	}
	decryptedFilePath, err := gpgPubkey.DecryptFile(tempFilePath)
	require.Error(t, err)
	require.Empty(t, decryptedFilePath)
}

func TestGPGDecryptEmptyPassphrase(t *testing.T) {
	tempFilePath := setupTestFile(t)
	gpgPubkey := GPG{
		PublicKey:  testPublicKey,
		PrivateKey: testPrivateKey,
		Passphrase: "",
	}

	encryptedFilePath, err := gpgPubkey.EncryptFile(tempFilePath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Remove(encryptedFilePath) })

	decryptedFilePath, err := gpgPubkey.DecryptFile(encryptedFilePath)
	require.Error(t, err)
	require.Empty(t, decryptedFilePath)
}

func TestGPGDecryptCorruptedEncryptedFile(t *testing.T) {
	tempFilePath := setupTestFile(t)
	gpgPubkey := GPG{
		PublicKey:  testPublicKey,
		PrivateKey: testPrivateKey,
		Passphrase: testPassphrase,
	}

	encryptedFilePath, err := gpgPubkey.EncryptFile(tempFilePath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Remove(encryptedFilePath) })

	f, err := os.OpenFile(encryptedFilePath, os.O_WRONLY, 0)
	require.NoError(t, err)
	_, err = f.WriteAt([]byte("corrupt"), 0)
	require.NoError(t, f.Close())
	require.NoError(t, err)

	decryptedFilePath, err := gpgPubkey.DecryptFile(encryptedFilePath)
	require.Error(t, err)
	require.Empty(t, decryptedFilePath)
}

func TestGPGDecryptFilePermissions(t *testing.T) {
	tempFilePath := setupTestFile(t)
	gpgPubkey := GPG{
		PublicKey:  testPublicKey,
		PrivateKey: testPrivateKey,
		Passphrase: testPassphrase,
	}

	encryptedFilePath, err := gpgPubkey.EncryptFile(tempFilePath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Remove(encryptedFilePath) })

	decryptedFilePath, err := gpgPubkey.DecryptFile(encryptedFilePath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Remove(decryptedFilePath) })

	info, err := os.Stat(decryptedFilePath)
	require.NoError(t, err)
	mode := info.Mode().Perm()
	require.Equal(t, os.FileMode(0600), mode)
}
