//go:build test_sqlize_repo
// +build test_sqlize_repo

package sqlize

import (
	"context"
	"database/sql"
	"net/url"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// Config is the required setting to open database connection.
type Config struct {
	Name       string
	Host       string
	User       string
	Pass       string
	Timezone   string
	SSLEnabled bool
}

// DSN returns config as DSN URI from.
func (c *Config) DSN() string {
	timezone := "utc"
	if len(c.Timezone) != 0 {
		timezone = c.Timezone
	}

	ssl := "disable"
	if c.SSLEnabled {
		ssl = "required"
	}

	q := make(url.Values)
	q.Set("timezone", timezone)
	q.Set("sslmode", ssl)

	dsn := url.URL{
		Scheme:   "postgres",
		Host:     c.Host,
		Path:     c.Name,
		User:     url.UserPassword(c.User, c.Pass),
		RawQuery: q.Encode(),
	}

	return dsn.String()
}

func SetupDBConnection() (*sql.DB, func(), error) {
	cfg := Config{
		Name:       os.Getenv("SQLIZE_DB_NAME"),
		Host:       os.Getenv("SQLIZE_DB_HOST"),
		User:       os.Getenv("SQLIZE_DB_USER"),
		Pass:       os.Getenv("SQLIZE_DB_PASS"),
		Timezone:   "Asia/Jakarta",
		SSLEnabled: false,
	}

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, func() {}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, func() {}, err
	}

	teardown := func() {
		_ = db.Close()
	}

	return db, teardown, nil
}

func TestPostgreRepository(t *testing.T) {
	db, teardown, err := SetupDBConnection()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "drop table if exists migrations_example;")
		teardown()
	})

	repo := NewPostgreRepository(db, "migrations_example")

	if err := repo.CreateMigrationTable(context.Background()); err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}

	migrations := []Migration{
		{
			Version: 1,
			Script:  "create table if not exists table_example_1(version bigint);",
		},
		{
			Script:  "insert into table_example_1(version) values(123);",
			Version: 2,
		},
		{
			Script:  "create table if not exists table_example_2(version bigint);",
			Version: 3,
		},
	}

	for _, m := range migrations {
		if err := repo.ApplyNewMigration(context.Background(), m.Version, m.Script, time.Now().Local()); err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}
	}

	histories, err := repo.FetchCurrentMigrations(context.Background())
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}

	if len(histories) != len(migrations) {
		t.Fatalf("length must be equals")
	}

	for _, m := range migrations {
		h, exist := histories[m.Version]
		if !exist {
			t.Fatalf("exepecting migration version %d exist in history", m.Version)
		}

		if !h.Applied {
			t.Fatalf("exepecting migration version %d exist is applied", m.Version)
		}
	}

	for i := len(migrations) - 1; i >= 0; i-- {
		m := migrations[i]
		if err := repo.UndoExistingMigration(context.Background(), m.Version, m.Script, time.Now().Local()); err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}
	}

	histories, err = repo.FetchCurrentMigrations(context.Background())
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}

	if len(histories) != len(migrations) {
		t.Fatalf("length must be equals")
	}

	for _, m := range migrations {
		h, exist := histories[m.Version]
		if !exist {
			t.Fatalf("exepecting migration version %d exist in history", m.Version)
		}

		if h.Applied {
			t.Fatalf("exepecting migration version %d exist is pending", m.Version)
		}
	}

	for _, m := range migrations {
		if err := repo.ApplyExistingMigration(context.Background(), m.Version, m.Script, time.Now().Local()); err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}
	}

	histories, err = repo.FetchCurrentMigrations(context.Background())
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}

	if len(histories) != len(migrations) {
		t.Fatalf("length must be equals")
	}

	for _, m := range migrations {
		h, exist := histories[m.Version]
		if !exist {
			t.Fatalf("exepecting migration version %d exist in history", m.Version)
		}

		if !h.Applied {
			t.Fatalf("exepecting migration version %d exist is applied", m.Version)
		}
	}

	for i := len(migrations) - 1; i >= 0; i-- {
		m := migrations[i]
		if err := repo.UndoExistingMigration(context.Background(), m.Version, m.Script, time.Now().Local()); err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}
	}

}
