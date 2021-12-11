package pqx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"
	"time"

	_ "github.com/lib/pq"
)

const (
	// Driver is a driver name.
	Driver = "postgres"
)

// Config is the required setting to open database connection.
type Config struct {
	Name              string
	Host              string
	User              string
	Pass              string
	Timezone          string
	SSLEnabled        bool
	MaxOpenConnection int
	MaxIdleConnection int
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
		Scheme:   Driver,
		Host:     c.Host,
		Path:     c.Name,
		User:     url.UserPassword(c.User, c.Pass),
		RawQuery: q.Encode(),
	}

	return dsn.String()
}

// Open knows how to open a new database connection.
// The Open just validate its arguments without creating a connection to the database.
// To verify that the data source name is valid, call database.CheckConnection.
//
// The returned DB is safe for concurrent use by multiple goroutines and maintains its own pool of idle connections.
// Thus, the Open function should be called just once.
// It is rarely necessary to close a DB.
func Open(cfg *Config) (*sql.DB, error) {
	dsn := cfg.DSN()
	db, err := sql.Open(Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: open db connection", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConnection)
	db.SetMaxIdleConns(cfg.MaxIdleConnection)

	return db, nil
}

// CheckConnection returns an error if connection is not ready.
// Otherwise, return nil.
// CheckConnection calls db.Ping to check if the connection is open,
// but that doesn't guarantee the database is ready to execute a query.
// So in addition to using db.Ping, CheckConnection also performs one-round-trip queries to the database.
//
// Use context.WithTimeout make a deadline.
func CheckConnection(ctx context.Context, maxTries int, db *sql.DB) error {
	pingErr := driver.ErrBadConn
	for tries := 0; pingErr != nil && tries < maxTries; tries++ {
		pingErr = db.PingContext(ctx)
		time.Sleep(time.Duration(tries) * 100 * time.Millisecond)
		// Canceled by deadline
		if ctx.Err() != nil {
			break
		}
	}

	// Make sure no context error occurs.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Make one round trip to database to make sure if the database ready to
	// handle query.
	_, err := db.ExecContext(ctx, "SELECT true;")
	if err != nil {
		return err
	}

	return nil
}
