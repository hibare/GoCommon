// Package hash provides utilities for hashing data.
package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// Hasher interface defines the methods for hashing data.
type Hasher interface {
	HashString(data string) (string, error)
	VerifyString(data string, hash string) (bool, error)
	HashFile(filePath string) (string, error)
	VerifyFile(filePath string, hash string) (bool, error)
}

// SHA256Hasher implements the Hasher interface for SHA-256 hashing.
type SHA256Hasher struct{}

// HashString hashes the given data and returns the hash as a hex string.
func (h *SHA256Hasher) HashString(data string) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write([]byte(data)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// HashFile hashes the given file and returns the hash as a hex string.
func (h *SHA256Hasher) HashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// VerifyString verifies if the hash of the given data matches the provided hash.
func (h *SHA256Hasher) VerifyString(data string, hash string) (bool, error) {
	calculatedHash, err := h.HashString(data)
	if err != nil {
		return false, err
	}
	return calculatedHash == hash, nil
}

// VerifyFile verifies if the hash of the given file matches the provided hash.
func (h *SHA256Hasher) VerifyFile(filePath string, hash string) (bool, error) {
	calculatedHash, err := h.HashFile(filePath)
	if err != nil {
		return false, err
	}
	return calculatedHash == hash, nil
}

func newSHA256Hasher() Hasher {
	return &SHA256Hasher{}
}

// NewSHA256Hasher returns a new SHA256Hasher.
var NewSHA256Hasher = newSHA256Hasher
