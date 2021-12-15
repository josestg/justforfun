package sqlize

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestFileSource_AppendMigration(t *testing.T) {

	t.Run("success scenario", func(t *testing.T) {
		temp, err := os.MkdirTemp("", "append-migration-*")
		if err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}
		t.Cleanup(func() {
			_ = os.RemoveAll(temp)
		})

		source := NewSourceFromDir(temp)

		path, err := source.AppendMigration(111, "example", DefaultTemplating.Template())
		if err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}

		expectedPath := filepath.Join(temp, "111_example")
		if path != expectedPath {
			t.Fatalf("expecting filepath %v but got %v", expectedPath, path)
		}
	})

	t.Run("destination file already taken", func(t *testing.T) {
		temp, err := os.MkdirTemp("", "append-migration-*")
		if err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}
		t.Cleanup(func() {
			_ = os.RemoveAll(temp)
		})

		source := NewSourceFromDir(temp)
		expectedPath := filepath.Join(temp, "112_example")

		_, err = os.OpenFile(expectedPath, os.O_CREATE|os.O_RDONLY, os.ModeExclusive)
		if err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}

		t.Cleanup(func() {
			os.Remove(expectedPath)
		})

		_, err = source.AppendMigration(112, "example", DefaultTemplating.Template())
		if !errors.Is(err, os.ErrPermission) {
			t.Fatalf("expecting error %v but got %v", os.ErrPermission, err)
		}

	})

	t.Run("dir is not exist", func(t *testing.T) {
		temp := fmt.Sprintf("%d", time.Now().UnixNano())

		source := NewSourceFromDir(temp)

		_, err := source.AppendMigration(111, "example", DefaultTemplating.Template())
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expecting error %v but got %v", os.ErrExist, err)
		}

	})
}

func TestFileSource_FetchMigrations(t *testing.T) {
	tmp, err := os.MkdirTemp("", "fetch-migration-*")
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(tmp)
	})

	files := []string{
		filepath.Join(tmp, "1_example_1.sql"),
		filepath.Join(tmp, "2_example_2.sql"),
		filepath.Join(tmp, "3_example_3.txt"),
		filepath.Join(tmp, "4_example_4.exe"),
		filepath.Join(tmp, "5_example_5.sql"),
		filepath.Join(tmp, "6_example_6.sql"),
	}

	for _, file := range files {
		func() {
			file, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				t.Fatalf("expecting error nil but got %v", err)
			}
			defer file.Close()

			io.WriteString(file, DefaultTemplating.Template())
		}()
	}

	// create some dir
	err = os.Mkdir(filepath.Join(tmp, "example"), os.ModeDir)
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}

	err = os.Mkdir(filepath.Join(tmp, "example.sql"), os.ModeDir)
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}

	// test part
	const upScriptExpected = `-- up script here...`
	const downScriptExpected = `-- down script here...`

	source := NewSourceFromDir(tmp)

	t.Run("success on migration up ascending", func(t *testing.T) {
		migrations, err := source.FetchMigrations(DefaultTemplating, MigrationUp, true)
		if err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}

		expectedUpMigration := []Migration{
			{
				Script:     upScriptExpected,
				Version:    1,
				SourcePath: "1_example_1.sql",
			},
			{
				Script:     upScriptExpected,
				Version:    2,
				SourcePath: "2_example_2.sql",
			},
			{
				Script:     upScriptExpected,
				Version:    5,
				SourcePath: "5_example_5.sql",
			},
			{
				Script:     upScriptExpected,
				Version:    6,
				SourcePath: "6_example_6.sql",
			},
		}

		if !reflect.DeepEqual(migrations, expectedUpMigration) {
			t.Fatalf("expecting equal")
		}
	})

	t.Run("success on migration up descending", func(t *testing.T) {
		migrations, err := source.FetchMigrations(DefaultTemplating, MigrationUp, false)
		if err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}

		expectedUpMigration := []Migration{
			{
				Script:     upScriptExpected,
				Version:    6,
				SourcePath: "6_example_6.sql",
			},
			{
				Script:     upScriptExpected,
				Version:    5,
				SourcePath: "5_example_5.sql",
			},
			{
				Script:     upScriptExpected,
				Version:    2,
				SourcePath: "2_example_2.sql",
			},
			{
				Script:     upScriptExpected,
				Version:    1,
				SourcePath: "1_example_1.sql",
			},
		}

		if !reflect.DeepEqual(migrations, expectedUpMigration) {
			t.Fatalf("expecting equal")
		}
	})

	t.Run("success on migration down", func(t *testing.T) {
		migrations, err := source.FetchMigrations(DefaultTemplating, MigrationDown, true)
		if err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}

		expectedUpMigration := []Migration{
			{
				Script:     downScriptExpected,
				Version:    1,
				SourcePath: "1_example_1.sql",
			},
			{
				Script:     downScriptExpected,
				Version:    2,
				SourcePath: "2_example_2.sql",
			},
			{
				Script:     downScriptExpected,
				Version:    5,
				SourcePath: "5_example_5.sql",
			},
			{
				Script:     downScriptExpected,
				Version:    6,
				SourcePath: "6_example_6.sql",
			},
		}

		if !reflect.DeepEqual(migrations, expectedUpMigration) {
			t.Fatalf("expecting equal")
		}
	})

	t.Run("success on migration status", func(t *testing.T) {
		migrations, err := source.FetchMigrations(DefaultTemplating, MigrationStatus, true)
		if err != nil {
			t.Fatalf("expecting error nil but got %v", err)
		}

		expectedUpMigration := []Migration{
			{
				Version:    1,
				SourcePath: "1_example_1.sql",
			},
			{
				Version:    2,
				SourcePath: "2_example_2.sql",
			},
			{
				Version:    5,
				SourcePath: "5_example_5.sql",
			},
			{
				Version:    6,
				SourcePath: "6_example_6.sql",
			},
		}

		if !reflect.DeepEqual(migrations, expectedUpMigration) {
			t.Fatalf("expecting equal")
		}
	})

}

func TestFileSource_FetchMigrations_MissingVersionFromFilename(t *testing.T) {
	tmp, err := os.MkdirTemp("", "fetch-migration-*")
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(tmp)
	})

	files := []string{
		filepath.Join(tmp, "1example1.sql"),
		filepath.Join(tmp, "2_example_2.sql"),
	}

	for _, file := range files {
		func() {
			file, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				t.Fatalf("expecting error nil but got %v", err)
			}
			defer file.Close()

			io.WriteString(file, DefaultTemplating.Template())
		}()
	}

	source := NewSourceFromDir(tmp)
	_, err = source.FetchMigrations(DefaultTemplating, MigrationUp, true)
	if !errors.Is(err, ErrVersionMissing) {
		t.Fatalf("expecting error %v but got %v", ErrVersionMissing, err)
	}

}

func TestFileSource_FetchMigrations_InvalidVersionType(t *testing.T) {
	tmp, err := os.MkdirTemp("", "fetch-migration-*")
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(tmp)
	})

	files := []string{
		filepath.Join(tmp, "1_example_1.sql"),
		filepath.Join(tmp, "-2_example_2.sql"),
	}

	for _, file := range files {
		func() {
			file, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				t.Fatalf("expecting error nil but got %v", err)
			}
			defer file.Close()

			io.WriteString(file, DefaultTemplating.Template())
		}()
	}

	source := NewSourceFromDir(tmp)
	_, err = source.FetchMigrations(DefaultTemplating, MigrationUp, true)
	if !errors.Is(err, ErrVersionTypeInvalid) {
		t.Fatalf("expecting error %v but got %v", ErrVersionTypeInvalid, err)
	}

}
