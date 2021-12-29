package iam

import (
	"context"
	"crypto"
	"crypto/rsa"

	"github.com/josestg/justforfun/pkg/xerrs"

	"github.com/josestg/justforfun/pkg/jwt"
)

var (
	ErrMissingSecretKeyName = xerrs.New("iam: missing secret key name")
)

type TokenEncoder interface {
	Encode(ctx context.Context, claims interface{}) (string, error)
}

type TokenDecoder interface {
	Decode(ctx context.Context, token string, claims interface{}) error
}

type Tokenizer interface {
	TokenDecoder
	TokenEncoder
}

type SecretKeyManager interface {
	SecretKey(ctx context.Context, keyName string) (secretKey interface{}, err error)
	NextSecretKey(ctx context.Context) (keyName string, secretKey interface{}, err error)
}

type JwtRS256 struct {
	alg string
	skm SecretKeyManager
}

func NewJwtRS256(skm SecretKeyManager) *JwtRS256 {
	return &JwtRS256{
		alg: "RS256",
		skm: skm,
	}
}

func (j *JwtRS256) Signer(ctx context.Context) (jwt.Signer, error) {
	keyName, secretKey, err := j.skm.NextSecretKey(ctx)
	if err != nil {
		return nil, xerrs.Wrap(err, "finding next secret key")
	}

	privateKey, compatible := secretKey.(*rsa.PrivateKey)
	if !compatible {
		return nil, xerrs.New("secret key must be type of *rsa.PrivateKey")
	}

	signer := jwt.NewRSASigner(keyName, crypto.SHA256, privateKey)
	return signer, nil
}

func (j *JwtRS256) Verifier(ctx context.Context) jwt.VerifierSelector {
	keySelector := jwt.VerifierSelector(func(header jwt.Header) (jwt.Verifier, error) {
		alg, exists := header["alg"].(string)
		if !exists || alg != j.alg {
			return nil, jwt.ErrInvalidFormat
		}

		kid, exists := header["kid"].(string)
		if !exists {
			return nil, ErrMissingSecretKeyName
		}

		secretKey, err := j.skm.SecretKey(ctx, kid)
		if err != nil {
			return nil, xerrs.Wrap(err, "finding corresponding secret key")
		}

		privateKey, compatible := secretKey.(*rsa.PrivateKey)
		if !compatible {
			return nil, xerrs.New("secret key must be type of *rsa.PrivateKey")
		}

		verifier := jwt.NewRSAVerifier(crypto.SHA256, &privateKey.PublicKey)
		return verifier, nil
	})

	return keySelector
}

type JwtSignVerifier interface {
	Verifier(ctx context.Context) jwt.VerifierSelector
	Signer(ctx context.Context) (jwt.Signer, error)
}

type JwtProvider struct {
	sv JwtSignVerifier
}

func NewJwtProvider(sv JwtSignVerifier) *JwtProvider {
	return &JwtProvider{
		sv: sv,
	}
}

func (j *JwtProvider) Decode(ctx context.Context, token string, claims interface{}) error {
	keySelector := j.sv.Verifier(ctx)
	if err := jwt.Decode(keySelector, token, claims); err != nil {
		return xerrs.Wrap(err, "decode token into claims")
	}

	return nil
}

func (j *JwtProvider) Encode(ctx context.Context, claims interface{}) (string, error) {
	jwtClaims, compatible := claims.(jwt.Valid)
	if !compatible {
		return "", xerrs.New("claims must be type of jwt.Valid")
	}

	signer, err := j.sv.Signer(ctx)
	if err != nil {
		return "", xerrs.Wrap(err, "creating singer")
	}

	header := jwt.Header{}
	token, err := jwt.Encode(signer, header, jwtClaims)
	if err != nil {
		return "", xerrs.Wrap(err, "creating jwt token")
	}

	return token, nil
}
