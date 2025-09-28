package gpg

import "errors"

var (
	// ErrKeyIDEmpty indicates a keyID cannot be empty
	ErrKeyIDEmpty = errors.New("keyID cannot be empty")

	// ErrKeyServerURLEmpty indicates a keyServerURL cannot be empty
	ErrKeyServerURLEmpty = errors.New("keyServerURL cannot be empty")

	// ErrEmptyInputFilePath indicates an inputFilePath cannot be empty
	ErrEmptyInputFilePath = errors.New("inputFilePath cannot be empty")

	// ErrNoEntitiesFoundInPrivateKey indicates no entities found in private key
	ErrNoEntitiesFoundInPrivateKey = errors.New("no entities found in private key")

	// ErrNoPrivateKeyFoundInEntity indicates no private key found in entity
	ErrNoPrivateKeyFoundInEntity = errors.New("no private key found in entity")
)
