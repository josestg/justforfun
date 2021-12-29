package user

import (
	"context"

	"github.com/josestg/justforfun/internal/usecase"
	"github.com/josestg/justforfun/pkg/validate"

	"github.com/josestg/justforfun/pkg/xerrs"

	"github.com/josestg/justforfun/internal/domain/user"
)

type Registration struct {
	p *usecase.Provider
}

func NewRegistration(provider *usecase.Provider) *Registration {
	return &Registration{
		p: provider,
	}
}

func (r *Registration) Register(ctx context.Context, input *user.RegistrationInput) (*user.User, error) {
	err := r.validateRegistrationInput(ctx, input)
	if err != nil {
		return nil, xerrs.Wrap(err, "validating registration input")
	}

	id := r.p.Identifier.NewUUID()
	u, err := user.CreateUser(id, input, r.p.Clock.Now())
	if err != nil {
		return nil, xerrs.Wrap(err, "creating user info")
	}

	err = r.p.Repository.User.Save(ctx, u)
	if err != nil {
		return nil, xerrs.Wrap(err, "save user info")
	}

	return u, nil
}

func (r *Registration) DeleteAccount(ctx context.Context, email, password string) error {
	u, err := r.p.Repository.User.FindByEmail(ctx, email)
	if err != nil {
		return xerrs.Wrap(err, "finding user by email")
	}

	if !u.ComparePasswordHash(password) {
		return user.ErrIncorrectCredential
	}

	err = r.p.Repository.User.Delete(ctx, u.ID, false)
	if err != nil {
		return xerrs.Wrap(err, "deleting user account")
	}

	return nil

}

func (r *Registration) validateRegistrationInput(ctx context.Context, input *user.RegistrationInput) error {
	schema := validate.Schema{
		"name":     validate.Field(&input.Name),
		"email":    validate.Field(&input.Email),
		"password": validate.Field(&input.Password),
	}

	return r.p.Validator.Validate(ctx, schema)
}
