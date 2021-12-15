package sqlize

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Migrator knows how to manage database migration.
type Migrator struct {
	printer    Printer
	repository Repository
	versioning VersionFunc
	templating TemplateReader
	source     Source
}

// NewMigrator creates a new Migrator instance.
func NewMigrator(source Source, repository Repository, options ...Option) *Migrator {
	m := Migrator{
		source:     source,
		repository: repository,
		printer:    DefaultPrinter,
		versioning: DefaultVersioning,
		templating: DefaultTemplating,
	}

	for _, fn := range options {
		fn(&m)
	}

	return &m
}

// Init knows how to init migrations.
func (m *Migrator) Init(ctx context.Context) error {
	if err := m.repository.CreateMigrationTable(ctx); err != nil {
		return fmt.Errorf("%w: creating migration table", err)
	}

	m.printer.Printf("migration table successfully initialized\n")
	return nil
}

// New creates a new migration source.
func (m *Migrator) New(name string) (string, error) {
	version := m.versioning()

	if !strings.HasPrefix(name, ".sql") {
		name = name + ".sql"
	}

	created, err := m.source.AppendMigration(version, name, m.templating.Template())
	if err != nil {
		return "", fmt.Errorf("%w: adding new migration into source", err)
	}

	m.printer.Printf("created new migration file at: %s", created)
	return created, nil
}

// Up knows how to apply all pending or created migrations.
func (m *Migrator) Up(ctx context.Context) error {
	histories, err := m.repository.FetchCurrentMigrations(ctx)
	if err != nil {
		return fmt.Errorf("%w: fetching migration histories", err)
	}

	migrations, err := m.source.FetchMigrations(m.templating, MigrationUp, true)
	if err != nil {
		return fmt.Errorf("%w: preparing migration", err)
	}

	latest := migrations[len(migrations)-1]
	if history, exist := histories[latest.Version]; exist && history.Applied {
		if len(migrations) != len(histories) {
			return errors.New("migrations history and source are mismatch")
		}

		m.printer.Printf("Already up to date.")
		return nil
	}

	appliedAt := time.Now().Local()
	for _, migration := range migrations {
		history, exist := histories[migration.Version]
		if !exist {
			err := m.repository.ApplyNewMigration(ctx, migration.Version, migration.Script, appliedAt)
			if err != nil {
				return fmt.Errorf("%w: applying new migration", err)
			}

			m.printer.Printf("%s   %s -> %s   %s\n", appliedAt.Format(time.Stamp), created, applied, migration.SourcePath)
			continue
		}

		if !history.Applied {
			err := m.repository.ApplyExistingMigration(ctx, migration.Version, migration.Script, appliedAt)
			if err != nil {
				return fmt.Errorf("%w: applying exsting migration", err)
			}

			m.printer.Printf("%s   %s -> %s   %s\n", appliedAt.Format(time.Stamp), pending, applied, migration.SourcePath)
		}
	}

	return nil
}

// Down knows how to reset migration into initial state.
func (m *Migrator) Down(ctx context.Context) error {
	histories, err := m.repository.FetchCurrentMigrations(ctx)
	if err != nil {
		return fmt.Errorf("%w: fetching migration histories", err)
	}

	if len(histories) == 0 {
		m.printer.Printf("Migration history is empty.")
		return nil
	}

	migrations, err := m.source.FetchMigrations(m.templating, MigrationDown, false)
	if err != nil {
		return fmt.Errorf("%w: preparing migration", err)
	}

	appliedAt := time.Now().Local()
	for _, migration := range migrations {
		history, exist := histories[migration.Version]
		if !exist {
			continue
		}

		if history.Applied {
			err := m.repository.UndoExistingMigration(ctx, migration.Version, migration.Script, appliedAt)
			if err != nil {
				return fmt.Errorf("%w: undo existing migration", err)
			}

			m.printer.Printf("%s   %s -> %s   %s\n", appliedAt.Format(time.Stamp), applied, pending, migration.SourcePath)
		}
	}

	return nil
}

// Undo knows how to undo one-step migration.
func (m *Migrator) Undo(ctx context.Context) error {
	histories, err := m.repository.FetchCurrentMigrations(ctx)
	if err != nil {
		return fmt.Errorf("%w: fetching migration histories", err)
	}

	if len(histories) == 0 {
		m.printer.Printf("Migration history is empty.\n")
		return nil
	}

	migrations, err := m.source.FetchMigrations(m.templating, MigrationDown, false)
	if err != nil {
		return fmt.Errorf("%w: preparing migration", err)
	}

	appliedAt := time.Now().Local()
	for _, migration := range migrations {
		history, exist := histories[migration.Version]
		if !exist {
			continue
		}

		if history.Applied {
			err := m.repository.UndoExistingMigration(ctx, migration.Version, migration.Script, appliedAt)
			if err != nil {
				return fmt.Errorf("%w: undo existing migration", err)
			}

			m.printer.Printf("%s   %s -> %s   %s\n", appliedAt.Format(time.Stamp), applied, pending, migration.SourcePath)
			break
		}
	}

	return nil
}

// Status prints the migration status.
func (m *Migrator) Status(ctx context.Context) error {
	histories, err := m.repository.FetchCurrentMigrations(ctx)
	if err != nil {
		return fmt.Errorf("%w: fetching migration histories", err)
	}

	migrations, err := m.source.FetchMigrations(m.templating, MigrationStatus, true)
	if err != nil {
		return fmt.Errorf("%w: preparing migration", err)
	}

	if len(migrations) == 0 {
		m.printer.Printf("migration source is empty\n")
		return nil
	}

	spaces := strings.Repeat("0", len(time.Stamp))
	for i := 0; i < len(migrations); i++ {
		history, exist := histories[migrations[i].Version]
		if !exist {
			m.printer.Printf("%s  %s -- %s\n", spaces, created, migrations[i].SourcePath)
			continue
		}

		if history.Applied {
			m.printer.Printf("%s  %s -- %s\n", history.DateApplied.Format(time.Stamp), applied, migrations[i].SourcePath)
		} else {
			m.printer.Printf("%s  %s -- %s\n", history.DateApplied.Format(time.Stamp), pending, migrations[i].SourcePath)
		}
	}

	return nil
}
