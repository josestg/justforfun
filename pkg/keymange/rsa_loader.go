package keymange

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

// Errors returned by RSA Signing Method and helpers
var (
	ErrKeyMustBePEMEncoded = errors.New("invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key")
	ErrNotRSAPrivateKey    = errors.New("key is not a valid RSA private key")
)

// ParseRSAPrivateKeyFromPEM is a helper method for
// parsing PEM encoded PKCS1 or PKCS8 private key
func ParseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}

// PemRSADir knows how to load PEM keys from dir and parses into rsa.PrivateKey.
func PemRSADir(dir fs.FS) ([]Key, error) {
	const ext = ".pem"
	var keys []Key

	err := fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.New("cannot walk into pemDir")
		}

		if d.IsDir() || filepath.Ext(path) != ext {
			return nil
		}

		file, err := dir.Open(path)
		if err != nil {
			return fmt.Errorf("open file %s: %w", path, err)
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("reading file %s: %w", path, err)
		}

		privateKey, err := ParseRSAPrivateKeyFromPEM(content)
		if err != nil {
			return fmt.Errorf("parsing pem: %w", err)
		}

		keys = append(keys, Key{
			Name:  strings.TrimSuffix(filepath.Base(path), ext),
			Value: privateKey,
		})

		return nil
	})

	return keys, err
}
