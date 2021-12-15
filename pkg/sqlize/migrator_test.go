//go:build test_sqlize_repo
// +build test_sqlize_repo

package sqlize

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func SetupMigrator() (*Migrator, *strings.Builder, func(), error) {
	db, dbTeardown, err := SetupDBConnection()
	if err != nil {
		return nil, nil, func() {}, err
	}

	tmp, err := os.MkdirTemp("", "testing-migration-*")
	if err != nil {
		return nil, nil, func() {}, nil
	}

	buf := &strings.Builder{}
	printer := PrinterFunc(func(format string, args ...interface{}) {
		buf.WriteString(fmt.Sprintf(format, args...))
	})

	migrationTable := strings.ReplaceAll(filepath.Base(tmp), "-", "_")

	repo := NewPostgreRepository(db, migrationTable)
	source := NewSourceFromDir(tmp)

	migrator := NewMigrator(
		source,
		repo,
		WithPrinter(printer),
		WithVersioning(DefaultVersioning),
		WithTemplating(DefaultTemplating),
	)

	teardown := func() {
		dbTeardown()
		_ = os.RemoveAll(tmp)
	}

	return migrator, buf, teardown, nil
}

func TestNewMigrator(t *testing.T) {
	migrator, logger, teardown, err := SetupMigrator()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(teardown)

	if err := migrator.Init(context.Background()); err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}

	initLog := "migration table successfully initialized"
	gotInitLog := logger.String()
	if gotInitLog != initLog {
		t.Fatalf("expecting init log %v but got %v", initLog, gotInitLog)
	}

	logger.Reset()

	if err := migrator.Status(context.Background()); err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}

	statusLog := "migration source is empty\n"
	gotStatusLog := logger.String()
	if gotStatusLog != statusLog {
		t.Fatalf("expecting init log %v but got %v", statusLog, gotStatusLog)
	}

	names := []string{"example_1", "example_2", "example_3"}
	migrationsFiles := make([]string, 0, len(names))

	for _, name := range names {
		path, err := migrator.New(name)
		if err != nil {
			t.Fatalf("expecting nil but got %v", err)
		}

		migrationsFiles = append(migrationsFiles, path)
	}

	if len(migrationsFiles) != len(names) {
		t.Fatalf("expecting length are equal")
	}

	scripts := []string{
		`
create table if not exists migrator_test_table_1(name text);
---+split+---
drop table if exists migrator_test_table_1;
`,
		`
insert into migrator_test_table_1(name) values('name_1');
insert into migrator_test_table_1(name) values('name_2');
---+split+---
delete from migrator_test_table_1 where name = 'name_1';
delete from migrator_test_table_1 where name = 'name_2';
`,
		`
alter table if exists migrator_test_table_1 add column qty int default 0;
---+split+---
alter table if exists migrator_test_table_1 drop column qty;
`,
	}

	for i := 0; i < len(migrationsFiles); i++ {
		func() {
			migration := migrationsFiles[i]
			script := scripts[i]
			file, err := os.OpenFile(migration, os.O_RDWR, os.ModePerm)
			if err != nil {
				t.Fatalf("expecting nil but got %v", err)
			}
			defer file.Close()

			_, err = file.WriteString(script)
			if err != nil {
				t.Fatalf("expecting nil but got %v", err)
			}
		}()
	}

	logger.Reset()
	if err := migrator.Status(context.Background()); err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}
	t.Log(logger.String())

	logger.Reset()
	if err := migrator.Up(context.Background()); err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}
	t.Log(logger.String())

	logger.Reset()
	if err := migrator.Status(context.Background()); err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}
	t.Log(logger.String())

	logger.Reset()
	if err := migrator.Down(context.Background()); err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}
	t.Log(logger.String())

	logger.Reset()
	if err := migrator.Up(context.Background()); err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}
	t.Log(logger.String())

	logger.Reset()
	if err := migrator.Undo(context.Background()); err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}

	if err := migrator.Undo(context.Background()); err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}
	t.Log(logger.String())

	logger.Reset()
	if err := migrator.Status(context.Background()); err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}
	t.Log(logger.String())

	logger.Reset()
	if err := migrator.Up(context.Background()); err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}
	t.Log(logger.String())

	logger.Reset()
	if err := migrator.Status(context.Background()); err != nil {
		t.Fatalf("expecting nil but got %v", err)
	}
	t.Log(logger.String())

}
