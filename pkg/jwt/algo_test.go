package jwt

import (
	"crypto"
	"crypto/rsa"
	"math/rand"
	"testing"
)

func TestAlgo_RSA(t *testing.T) {
	reader := rand.New(rand.NewSource(1))
	private, err := rsa.GenerateKey(reader, 2048)
	if err != nil {
		t.Errorf("expecting error nil but got %v", err)
	}

	const content = "content-example"
	kid := "kid-example"
	signer := NewRSASigner(kid, crypto.SHA256, private)

	signature, err := signer.Sign([]byte(content))
	if err != nil {
		t.Errorf("expecting error nil but got %v", err)
	}

	verifier := NewRSAVerifier(crypto.SHA256, &private.PublicKey)

	err = verifier.Verify([]byte(content), signature)
	if err != nil {
		t.Errorf("expecting error nil but got %v", err)
	}

	err = verifier.Verify([]byte("invalid content"), signature)
	if err == nil {
		t.Errorf("expecting error non-nil")
	}
}
