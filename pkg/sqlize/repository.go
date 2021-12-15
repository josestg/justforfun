package sqlize

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// postgre implements Repository for PostgreSQL Database.
type postgre struct {
	db    *sql.DB
	table string
}

// NewPostgreRepository creates a new Repository for PostgreSQL Database.
func NewPostgreRepository(db *sql.DB, table string) Repository {
	return &postgre{
		db:    db,
		table: table,
	}
}

func (p *postgre) CreateMigrationTable(ctx context.Context) error {
	query := fmt.Sprintf(`
create table if not exists %s(
    id serial not null,
    migration_version bigint not null,
    applied boolean default false,
    date_applied timestamp default current_date,
    
	constraint %s_primary_key primary key (id),
	constraint %s_unique_version unique(migration_version)
);
`, p.table, p.table, p.table)

	if _, err := p.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("%w: creating migration table: %s", err, p.table)
	}

	return nil
}

func (p *postgre) FetchCurrentMigrations(ctx context.Context) (map[int64]MigrationHistory, error) {
	query := fmt.Sprintf(`
select migration_version, applied, date_applied 
from %s order by migration_version;
`, p.table)

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: selecting current migrations", err)
	}
	defer rows.Close()

	histories := make(map[int64]MigrationHistory)

	for rows.Next() {
		var history MigrationHistory
		if err := rows.Scan(&history.Version, &history.Applied, &history.DateApplied); err != nil {
			return nil, fmt.Errorf("%w: scanning history record", err)
		}

		histories[history.Version] = history
	}

	return histories, nil
}

func (p *postgre) ApplyNewMigration(ctx context.Context, version int64, script string, appliedAt time.Time) error {
	return p.applyMigration(ctx, version, true, script, appliedAt)
}

func (p *postgre) ApplyExistingMigration(ctx context.Context, version int64, script string, appliedAt time.Time) error {
	return p.updateMigration(ctx, version, true, script, appliedAt)
}

func (p *postgre) UndoExistingMigration(ctx context.Context, version int64, script string, appliedAt time.Time) error {
	return p.updateMigration(ctx, version, false, script, appliedAt)
}

func (p *postgre) updateMigration(ctx context.Context, version int64, applied bool, script string, appliedAt time.Time) error {
	query := fmt.Sprintf(`update %s set applied = $2, date_applied = $3 where migration_version = $1;`, p.table)
	return p.transaction(ctx, script, query, version, applied, appliedAt)
}

func (p *postgre) applyMigration(ctx context.Context, version int64, applied bool, script string, appliedAt time.Time) error {
	query := fmt.Sprintf(`insert into %s (migration_version, applied, date_applied) values ($1, $2, $3);`, p.table)
	return p.transaction(ctx, script, query, version, applied, appliedAt)
}

func (p *postgre) transaction(ctx context.Context, script string, query string, version int64, applied bool, appliedAt time.Time) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("%w: creating  transaction", err)
	}

	if _, err := tx.ExecContext(ctx, script); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: executing script", err)
	}

	if _, err := tx.ExecContext(ctx, query, version, applied, appliedAt); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: inserting migration version: %d", err, version)
	}

	return tx.Commit()
}
