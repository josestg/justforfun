package sqlize

import (
	"context"
	"fmt"
	"io"
	"time"
)

// Action is a type for available actions.
type Action uint

const (
	MigrationUp     = Action(0)
	MigrationDown   = Action(1)
	MigrationStatus = Action(2)
)

const (
	pending = "pending"
	applied = "applied"
	created = "created"
)

// DefaultTemplating is a default templating engine.
// This templating engine use token `---+split+---` to split migration actions.
// for example:
//
//  create table if not exists migrator_test_table_1(name text);
// 	---+split+---
// 	drop table if exists migrator_test_table_1;
var DefaultTemplating = &templating{}

// DefaultVersioning is a default version generator using time unix nano.
var DefaultVersioning = func() int64 { return time.Now().UnixNano() }

// DefaultPrinter is a default logger for migrator.
var DefaultPrinter = PrinterFunc(func(format string, args ...interface{}) { fmt.Printf(format, args...) })

// VersionFunc is a signature for version generator.
type VersionFunc func() int64

// Option is an option type that can be used to customize the Migrator configs.
type Option func(m *Migrator)

// Migration represents a migration unit.
type Migration struct {
	Script     string
	Version    int64
	SourcePath string
}

// MigrationHistory represent a migration history unit.
type MigrationHistory struct {
	Migration
	Applied     bool
	DateApplied time.Time
}

// Repository knows how to manage migration at persistence storage level.
type Repository interface {
	// CreateMigrationTable knows how to create the migration table.
	CreateMigrationTable(ctx context.Context) error

	// FetchCurrentMigrations know how to fetch current migration histories.
	FetchCurrentMigrations(ctx context.Context) (map[int64]MigrationHistory, error)

	// ApplyNewMigration knows how to apply a new migration.
	ApplyNewMigration(ctx context.Context, version int64, script string, appliedAt time.Time) error

	// ApplyExistingMigration knows how to apply an existing migration.
	ApplyExistingMigration(ctx context.Context, version int64, script string, appliedAt time.Time) error

	// UndoExistingMigration knows how to undo existing migration.
	UndoExistingMigration(ctx context.Context, version int64, script string, appliedAt time.Time) error
}

// Source knows how to manage migration source.
// We can implement this contract to manage migration source from
// local file system (for example: FileSource)
type Source interface {
	// AppendMigration inserts a new migration into source.
	AppendMigration(version int64, name string, content string) (string, error)
	// FetchMigrations knows how to get list of migrations from source.
	FetchMigrations(reader Reader, action Action, ascending bool) ([]Migration, error)
}

// Template knows how to manage template.
type Template interface {
	// Template returns a template format.
	Template() string
}

// TemplateReader knows how to manage and read template.
type TemplateReader interface {
	Template
	Reader
}

// Reader knows how to read and parse template.
type Reader interface {
	// Read knows how to read and parse template based on given action.
	Read(reader io.Reader, action Action) (string, error)
}

// Printer is a contract for migration logger.
type Printer interface {
	// Printf knows how to print formatted text.
	Printf(format string, args ...interface{})
}

// PrinterFunc is adapter function that can be used to create a new Printer
// by using function signature.
type PrinterFunc func(format string, args ...interface{})

func (p PrinterFunc) Printf(format string, args ...interface{}) { p(format, args...) }
