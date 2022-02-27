package repository

import (
	"database/sql"

	"github.com/josestg/justforfun/internal/domain/user"

	"github.com/josestg/justforfun/internal/repository/usersql"
)

// Container contains all repository instances.
type Container struct {
	User user.Repository
}

// NewSQLContainer creates a new SQL-based repository container.
func NewSQLContainer(db *sql.DB) *Container {
	return &Container{
		User: usersql.NewRepository(db),
	}
}
