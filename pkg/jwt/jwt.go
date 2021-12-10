package jwt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrExpired       = errors.New("jwt: claims expired")
	ErrNotBefore     = errors.New("jwt: claims not active yet")
	ErrInvalidFormat = errors.New("jwt: invalid token format")
)

// Header represents JWT header.
type Header map[string]interface{}

// VerifierSelector selects a correct verifier based on given Header.
type VerifierSelector func(header Header) (Verifier, error)

// Valid knows how validate claims.
type Valid interface {
	// Valid returns an error if the claims is not valid at the given time.
	Valid(at *Time) error
}

// Signer is contract for signer algorithm.
// see: RSASigner.
type Signer interface {
	// Sign signs the given bytes with the implementation algorithm.
	Sign(b []byte) (signature []byte, err error)

	// Header returns a required headers to make Verifier understand
	// the signature.
	Header() Header
}

// Verifier is contract for signature verification.
type Verifier interface {
	// Verify verifies payload and signature using the implementation algorithm.
	Verify(payload, signature []byte) error
}

// Encode encodes the header and payload (claims) into a signed JWT.
func Encode(singer Signer, header Header, payload interface{}) (string, error) {
	// inject the signer's header into user's defined header.
	// we can add 'alg' and 'kid' headers here.
	for k, v := range singer.Header() {
		header[k] = v
	}

	headerEncoded, err := b64URLEncoded(header)
	if err != nil {
		return "", fmt.Errorf("%w: encode header part", err)
	}

	payloadEncoded, err := b64URLEncoded(payload)
	if err != nil {
		return "", fmt.Errorf("%w: encode payload part", err)
	}

	// JWT formats:
	// 	base64UrlEncode(header) + "." + base64UrlEncode(payload) + "." + base64UrlEncode(signature)
	sb := bytes.Buffer{}
	sb.Write(headerEncoded)
	sb.WriteByte('.')
	sb.Write(payloadEncoded)

	signature, err := singer.Sign(sb.Bytes())
	if err != nil {
		return "", fmt.Errorf("%w: create signature", err)
	}

	sb.WriteRune('.') // appends last '.' before the encoded signature.
	_, err = sb.WriteString(base64.RawURLEncoding.EncodeToString(signature))
	return sb.String(), err
}

// Decode decodes the given token and store the result into pointed claims.
// Before token decoded, the selector choose a proper algorithm to verify the token signature
// based on token's header.
func Decode(selector VerifierSelector, token string, claims interface{}) error {
	// JWT formats:
	// 	base64UrlEncode(header) + "." + base64UrlEncode(payload) + "." + base64UrlEncode(signature)
	parts := strings.SplitN(token, ".", 3)
	if len(parts) != 3 {
		return ErrInvalidFormat
	}

	var header Header
	if err := b64URLEncodeToJSON(parts[0], &header); err != nil {
		return fmt.Errorf("%w: encode base64-url header into Header", err)
	}

	verifier, err := selector(header)
	if err != nil {
		return fmt.Errorf("%w: selecting token verifier", err)
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return fmt.Errorf("%w: creating signature", err)
	}

	content := fmt.Sprintf("%s.%s", parts[0], parts[1])
	if err := verifier.Verify([]byte(content), signature); err != nil {
		return fmt.Errorf("%w: verifying signature", err)
	}

	if err := b64URLEncodeToJSON(parts[1], claims); err != nil {
		return fmt.Errorf("%w: encode base64-url body into claims", err)
	}

	if t, ok := claims.(Valid); ok {
		return t.Valid(NewTime(time.Now()))
	}

	return nil
}

func b64URLEncodeToJSON(urlEncoded string, v interface{}) error {
	reader := base64.NewDecoder(base64.RawURLEncoding, strings.NewReader(urlEncoded))
	return json.NewDecoder(reader).Decode(v)
}

func b64URLEncoded(v interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	encoder := base64.NewEncoder(base64.RawURLEncoding, &buf)
	if err := json.NewEncoder(encoder).Encode(v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
