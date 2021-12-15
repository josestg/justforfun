package sqlize

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrVersionMissing     = errors.New("sqlize: version part is missing")
	ErrVersionTypeInvalid = errors.New("sqlize: version type must be int64 and must be a positive number")
)

type FileSource struct {
	dir string
}

func NewSourceFromDir(dir string) Source {
	return &FileSource{
		dir: dir,
	}
}

func (f *FileSource) AppendMigration(version int64, name string, content string) (string, error) {
	fullName := fmt.Sprintf("%d_%s", version, name)
	fullPath := filepath.Join(f.dir, fullName)

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("%w: creating new migration file: %s", err, fullPath)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return "", fmt.Errorf("%w: writing migration template", err)
	}

	return fullPath, nil
}

func (f *FileSource) FetchMigrations(reader Reader, action Action, ascending bool) ([]Migration, error) {
	migrations := make([]Migration, 0)

	dir, err := os.Open(f.dir)
	if err != nil {
		return nil, fmt.Errorf("%w: checking dir", err)
	}
	defer dir.Close()

	fsSys := os.DirFS(f.dir)
	err = fs.WalkDir(fsSys, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || err != nil {
			return nil
		}

		filename := d.Name()
		if !strings.HasSuffix(filename, ".sql") {
			return nil
		}

		parts := strings.SplitN(filename, "_", 2)
		if len(parts) != 2 {
			return ErrVersionMissing
		}

		version, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return ErrVersionTypeInvalid
		}

		if version <= 0 {
			return ErrVersionTypeInvalid
		}

		migration := Migration{
			Version:    version,
			SourcePath: filename,
			Script:     "", // fill if needed.
		}

		switch action {
		case MigrationDown, MigrationUp:
			file, err := fsSys.Open(path)
			if err != nil {
				return fmt.Errorf("%w: open migration file: %s", err, filename)
			}

			script, err := reader.Read(file, action)
			if err != nil {
				return fmt.Errorf("%w: reading script from template", err)
			}

			migration.Script = script
		}

		migrations = append(migrations, migration)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: collecting migration from source dir", err)
	}

	sort.Slice(migrations, func(i, j int) bool {
		if ascending {
			return migrations[i].Version < migrations[j].Version
		}

		return migrations[j].Version < migrations[i].Version
	})

	return migrations, nil
}
