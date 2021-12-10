package jwt

import (
	"crypto"
	"crypto/rsa"
	"errors"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestSign_RSA(t *testing.T) {
	reader := rand.New(rand.NewSource(1))
	private, err := rsa.GenerateKey(reader, 2048)
	if err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}

	kid := "kid-example"
	signer := NewRSASigner(kid, crypto.SHA256, private)

	claims := StandardClaims{
		ID:        "123",
		Issuer:    "just-for-func",
		Subject:   "12345",
		IssuedAt:  NewTime(time.Now()),
		ExpiresAt: NewTime(time.Now().Add(time.Hour)),
	}

	token, err := Encode(signer, Header{}, claims)
	if err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}

	selector := func(header Header) (Verifier, error) {
		if header["alg"] == "RS256" {
			return &RSAVerifier{
				public: &private.PublicKey,
				hash:   crypto.SHA256,
			}, nil
		}
		return nil, errors.New("unknown alg")
	}

	var claims2 StandardClaims
	if err = Decode(selector, token, &claims2); err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}

	if !reflect.DeepEqual(claims, claims2) {
		t.Fatalf("expecting claims are equal")
	}
}

var fakeErr = errors.New("fake error")

type signerMock struct {
	err error
}

func (s *signerMock) Sign(b []byte) (signature []byte, err error) {
	if s.err != nil {
		return nil, s.err
	}

	return b, nil
}

func (s *signerMock) Header() Header {
	return Header{
		"kid": "kid-1",
		"alg": "alg-1",
	}
}

func TestEncode(t *testing.T) {
	t.Run("if signer not failed", func(t *testing.T) {
		signer := &signerMock{}

		header := Header{}
		claims := []byte("example")

		token, err := Encode(signer, header, claims)
		if err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}

		parts := strings.SplitN(token, ".", 3)
		if len(parts) != 3 {
			t.Fatalf("expecting contains 3 parts but got %v part(s)", len(parts))
		}

		var tokenHeader Header
		if err := b64URLEncodeToJSON(parts[0], &tokenHeader); err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}

		if !reflect.DeepEqual(tokenHeader, signer.Header()) {
			t.Fatalf("expecting header values are equals")
		}
	})

	t.Run("if signer failed", func(t *testing.T) {
		signer := &signerMock{err: fakeErr}

		header := Header{}
		claims := []byte("example")

		token, err := Encode(signer, header, claims)
		if err == nil {
			t.Fatalf("expecting error not nil")
		}

		if len(token) != 0 {
			t.Fatalf("execting token is empty")
		}

		if !errors.Is(err, fakeErr) {
			t.Fatalf("expecting error from signer")
		}
	})
}

type verifierMock struct {
	err error
}

func (v *verifierMock) Verify(payload, signature []byte) error {
	return v.err
}

func TestDecode(t *testing.T) {
	t.Run("verifier succeed", func(t *testing.T) {
		selector := func(header Header) (Verifier, error) {
			return &verifierMock{}, nil
		}

		const token = `eyJhbGciOiJhbGctMSIsImtpZCI6ImtpZC0xIn0K.IlpYaGhiWEJzWlE9PSIK.ZXlKaGJHY2lPaUpoYkdjdE1TSXNJbXRwWkNJNkltdHBaQzB4SW4wSy5JbHBZYUdoaVdFSnpXbEU5UFNJSw`

		var claims string
		if err := Decode(selector, token, &claims); err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}

		if len(claims) == 0 {
			t.Fatalf("expecting claims not empty")
		}
	})

	t.Run("invalid token format", func(t *testing.T) {
		selector := func(header Header) (Verifier, error) {
			return &verifierMock{}, nil
		}

		const token = `IlpYaGhiWEJzWlE9PSIK.ZXlKaGJHY2lPaUpoYkdjdE1TSXNJbXRwWkNJNkltdHBaQzB4SW4wSy5JbHBZYUdoaVdFSnpXbEU5UFNJSw`

		var claims string
		if err := Decode(selector, token, &claims); err != ErrInvalidFormat {
			t.Fatalf("expecting error %v but got %v", ErrInvalidFormat, err)
		}
	})

	t.Run("header is not valid form", func(t *testing.T) {
		selector := func(header Header) (Verifier, error) {
			return &verifierMock{}, nil
		}

		// invalid header => header is not a valid jwt.Header type.
		const token = `IlpYaGhiWEJzWlE9PSIK.IlpYaGhiWEJzWlE9PSIK.ZXlKaGJHY2lPaUpoYkdjdE1TSXNJbXRwWkNJNkltdHBaQzB4SW4wSy5JbHBZYUdoaVdFSnpXbEU5UFNJSw`

		var claims string
		if err := Decode(selector, token, &claims); err == nil {
			t.Fatalf("expecting error not-nil")
		}

		if len(claims) != 0 {
			t.Fatalf("expecting claims empty")
		}
	})

	t.Run("payload is not valid form", func(t *testing.T) {
		selector := func(header Header) (Verifier, error) {
			return &verifierMock{}, nil
		}

		// body is not a valid json format.
		const token = `eyJhbGciOiJhbGctMSIsImtpZCI6ImtpZC0xIn0K.ZXlKaGJHY2lPaUpoYkdjdE1TSXNJbXRwWkNJNkltdHBaQzB4SW4wSy5JbHBZYUdoaVdFSnpXbEU5UFNJSw.ZXlKaGJHY2lPaUpoYkdjdE1TSXNJbXRwWkNJNkltdHBaQzB4SW4wSy5JbHBZYUdoaVdFSnpXbEU5UFNJSw`

		var claims string
		if err := Decode(selector, token, &claims); err == nil {
			t.Fatalf("expecting error not-nil")
		}

		if len(claims) != 0 {
			t.Fatalf("expecting claims empty")
		}
	})

	t.Run("selector return error", func(t *testing.T) {
		selector := func(header Header) (Verifier, error) {
			return nil, fakeErr
		}

		// body is not a valid json format.
		const token = `eyJhbGciOiJhbGctMSIsImtpZCI6ImtpZC0xIn0K.ZXlKaGJHY2lPaUpoYkdjdE1TSXNJbXRwWkNJNkltdHBaQzB4SW4wSy5JbHBZYUdoaVdFSnpXbEU5UFNJSw.ZXlKaGJHY2lPaUpoYkdjdE1TSXNJbXRwWkNJNkltdHBaQzB4SW4wSy5JbHBZYUdoaVdFSnpXbEU5UFNJSw`

		var claims string
		if err := Decode(selector, token, &claims); !errors.Is(err, fakeErr) {
			t.Fatalf("expecting error is a returned error from selector")
		}

		if len(claims) != 0 {
			t.Fatalf("expecting claims empty")
		}
	})

	t.Run("verifier failed", func(t *testing.T) {
		selector := func(header Header) (Verifier, error) {
			return &verifierMock{err: fakeErr}, nil
		}

		// body is not a valid json format.
		const token = `eyJhbGciOiJhbGctMSIsImtpZCI6ImtpZC0xIn0K.ZXlKaGJHY2lPaUpoYkdjdE1TSXNJbXRwWkNJNkltdHBaQzB4SW4wSy5JbHBZYUdoaVdFSnpXbEU5UFNJSw.ZXlKaGJHY2lPaUpoYkdjdE1TSXNJbXRwWkNJNkltdHBaQzB4SW4wSy5JbHBZYUdoaVdFSnpXbEU5UFNJSw`

		var claims string
		if err := Decode(selector, token, &claims); !errors.Is(err, fakeErr) {
			t.Fatalf("expecting error is a returned error from verifier")
		}

		if len(claims) != 0 {
			t.Fatalf("expecting claims empty")
		}
	})
}
