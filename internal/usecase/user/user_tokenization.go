package user

import (
	"context"
	"time"

	"github.com/josestg/justforfun/internal/usecase"

	"github.com/josestg/justforfun/pkg/xerrs"

	"github.com/josestg/justforfun/internal/domain/user"

	"github.com/josestg/justforfun/internal/iam"
)

type Tokenization struct {
	tokenizer       iam.Tokenizer
	tokenExpiration time.Duration
	p               *usecase.Provider
}

func NewTokenization(provider *usecase.Provider, tokenizer iam.Tokenizer) *Tokenization {
	return &Tokenization{
		tokenExpiration: time.Hour,
		tokenizer:       tokenizer,
		p:               provider,
	}
}

func (t *Tokenization) ParseToken(ctx context.Context, token string) (*user.JwtClaims, error) {
	var claims user.JwtClaims
	if err := t.tokenizer.Decode(ctx, token, &claims); err != nil {
		return nil, xerrs.Wrap(err, "decode claims on parsing token")
	}

	return &claims, nil
}

func (t *Tokenization) GenerateToken(ctx context.Context, email string, password string) (string, error) {
	u, err := t.p.Repository.User.FindByEmail(ctx, email)
	if err != nil {
		return "", xerrs.Wrap(err, "finding user by email")
	}

	if !u.ComparePasswordHash(password) {
		return "", user.ErrIncorrectCredential
	}

	claims := user.NewJwtClaims(u, t.tokenExpiration)
	token, err := t.tokenizer.Encode(ctx, claims)
	if err != nil {
		return "", xerrs.Wrap(err, "encode claims on generating token")
	}

	return token, nil
}
