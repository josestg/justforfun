package usersql

import (
	"context"
	"database/sql"

	"github.com/josestg/justforfun/pkg/xerrs"

	"github.com/josestg/justforfun/internal/domain/user"
)

// Repository implement user repository interface.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new user repository.
func NewRepository(db *sql.DB) user.Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(ctx context.Context, u *user.User) error {
	const q = `insert into users(id, name, email, password_hash, date_created, date_updated) values ($1, $2, $3, $4, $5, $6);`

	_, err := r.db.ExecContext(ctx, q, u.ID, u.Name, u.Email, u.HashedPassword, u.DateCreated, u.DateUpdated)
	if err != nil {
		return xerrs.Wrap(err, "executing insert query.")
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string, hard bool) error {
	if hard {
		return r.hardDelete(ctx, id)
	}

	return r.softDelete(ctx, id)
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	const q = `select id, name, email, password_hash, date_created, date_updated from users where email = $1`
	row := r.db.QueryRowContext(ctx, q, email)

	var u user.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.DateCreated, &u.DateUpdated); err != nil {
		return nil, xerrs.Wrap(err, "querying user by email")
	}

	return &u, nil
}

func (r *Repository) softDelete(ctx context.Context, id string) error {

	return nil
}

func (r *Repository) hardDelete(ctx context.Context, id string) error {
	return nil
}
