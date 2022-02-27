package user

import (
	"context"
	"time"

	"github.com/josestg/justforfun/pkg/jwt"

	"golang.org/x/crypto/bcrypt"

	"github.com/josestg/justforfun/pkg/xerrs"
)

var (
	ErrIncorrectCredential = xerrs.New("user.incorrect_credential")
)

type Repository interface {
	Save(ctx context.Context, u *User) error
	Delete(ctx context.Context, id string, hard bool) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type JwtClaims struct {
	jwt.StandardClaims
	UserProfileName string `json:"user_profile_name"`
}

func NewJwtClaims(u *User, exp time.Duration) *JwtClaims {
	now := time.Now()
	return &JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   u.ID,
			IssuedAt:  jwt.NewTime(now),
			ExpiresAt: jwt.NewTime(now.Add(exp)),
		},
		UserProfileName: u.Name,
	}
}

// User contains information about user.
type User struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	HashedPassword string     `json:"-"`
	DateCreated    time.Time  `json:"date_created"`
	DateUpdated    time.Time  `json:"date_updated"`
	DateDeleted    *time.Time `json:"date_deleted"`
}

func CreateUser(id string, input *RegistrationInput, timestamp time.Time) (*User, error) {
	hashed, err := HashPassword(input.Password)
	if err != nil {
		return nil, xerrs.Wrap(err, "hashing password")
	}

	u := &User{
		ID:             id,
		Name:           input.Name,
		Email:          input.Email,
		HashedPassword: hashed,
		DateCreated:    timestamp,
		DateUpdated:    timestamp,
	}

	return u, nil
}

func (u *User) ComparePasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", xerrs.Wrap(err, "generating password hash")
	}

	return string(b), nil
}

type RegistrationInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
