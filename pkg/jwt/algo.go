package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

// RSAVerifier knows how to verify payload and signature using
// the RSA algorithm.
type RSAVerifier struct {
	public *rsa.PublicKey
	hash   crypto.Hash
}

// NewRSAVerifier creates a new verifier using RSA public key.
func NewRSAVerifier(hash crypto.Hash, public *rsa.PublicKey) *RSAVerifier {
	return &RSAVerifier{
		hash:   hash,
		public: public,
	}
}

func (r *RSAVerifier) Verify(payload, signature []byte) error {
	hashed := sha256.Sum256(payload)

	if err := rsa.VerifyPKCS1v15(r.public, r.hash, hashed[:], signature); err != nil {
		return fmt.Errorf("%w: verifying signature using RSA", err)
	}

	return nil
}

// RSASigner knows how to sign a given payload using the RSA algorithm.
type RSASigner struct {
	header  Header
	private *rsa.PrivateKey
	hash    crypto.Hash
}

// NewRSASigner creates a new RSA signer.
func NewRSASigner(kid string, hash crypto.Hash, private *rsa.PrivateKey) *RSASigner {
	return &RSASigner{
		header: Header{
			"kid": kid,
			"alg": "RS256",
		},
		private: private,
		hash:    hash,
	}
}

func (r *RSASigner) Sign(b []byte) ([]byte, error) {
	hashed := sha256.Sum256(b)
	
	signature, err := rsa.SignPKCS1v15(rand.Reader, r.private, r.hash, hashed[:])
	if err != nil {
		return nil, fmt.Errorf("%w: signing using RSA", err)
	}

	return signature, nil
}

func (r *RSASigner) Header() Header {
	return r.header
}
