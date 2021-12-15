package conf

import (
	"time"

	"github.com/josestg/justforfun/pkg/env"
	"github.com/josestg/justforfun/pkg/pqx"
)

// Option is option type for customize the Config.
type Option func(c *Config)

// Config holds all configs.
type Config struct {
	DB        *DB        `json:"db,omitempty"`
	RestAPI   *RestAPI   `json:"rest_api,omitempty"`
	Migration *Migration `json:"migration"`
}

// New creates a new config based on given options.
func New(options ...Option) *Config {
	c := Config{
		DB:        &DB{},
		RestAPI:   &RestAPI{},
		Migration: &Migration{},
	}

	for _, fn := range options {
		fn(&c)
	}

	return &c
}

// RestAPI hold all REST API config.
type RestAPI struct {
	Addr            string        `json:"addr"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
}

// WithRestAPIFromOSEnv creates a RestAPI config loader from OS Env.
func WithRestAPIFromOSEnv() Option {
	return func(c *Config) {
		c.RestAPI = &RestAPI{
			Addr:            env.String("API_ADDR", ":8000"),
			ReadTimeout:     env.Duration("API_REQUEST_READ_TIMEOUT", 20*time.Second),
			WriteTimeout:    env.Duration("API_REQUEST_WRITE_TIMEOUT", 30*time.Second),
			ShutdownTimeout: env.Duration("API_REQUEST_SHUTDOWN_TIMEOUT", 30*time.Second),
		}
	}
}

// Migration holds all Migration config.
type Migration struct {
	SourceDir string `json:"source_dir"`
	TableName string `json:"table_name"`
}

// WithMigrationFromEnv creates a Migration config loader from OS Env.
func WithMigrationFromEnv() Option {
	return func(c *Config) {
		c.Migration = &Migration{
			SourceDir: env.String("MIGRATION_SOURCE_DIR", "vars/migrations"),
			TableName: env.String("MIGRATION_TABLE_NAME", "sqlize_migrations"),
		}
	}
}

// DB holds all DB config.
type DB struct {
	Postgre *pqx.Config `json:"postgre,omitempty"`
}

// WithDBPostgreFromOSEnv creates a DB config loader from OS Env.
func WithDBPostgreFromOSEnv() Option {
	return func(c *Config) {
		c.DB.Postgre = &pqx.Config{
			Name:              env.String("POSTGRE_DB_NAME", "justforfun"),
			Host:              env.String("POSTGRE_DB_HOST", "0.0.0.0:5432"),
			User:              env.String("POSTGRE_DB_USER", "postgres"),
			Pass:              env.String("POSTGRE_DB_PASS", "kunci"),
			Timezone:          env.String("POSTGRE_DB_TIMEZONE", "Asia/Jakarta"),
			SSLEnabled:        env.Bool("POSTGRE_DB_SSL_ENABLED", false),
			MaxOpenConnection: env.Int("POSTGRE_DB_MAX_OPEN_CONN", 8),
			MaxIdleConnection: env.Int("POSTGRE_DB_MAX_IDLE_CONN", 8),
		}
	}
}
