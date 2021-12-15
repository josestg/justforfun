package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/josestg/justforfun/pkg/pqx"

	"github.com/josestg/justforfun/pkg/sqlize"

	"github.com/josestg/justforfun/pkg/xerrs"
)

func main() {
	args := os.Args[1:]
	if err := run(args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		args = append(args, "help")
	}

	// open database connection.
	//
	// note: we only open connection once,
	// if we need database connection we must pass it as dependency.
	db, err := pqx.Open(&pqx.Config{
		Name:              "postgres",
		Host:              "localhost:5432",
		User:              "postgres",
		Pass:              "kunci",
		Timezone:          "Asia/Jakarta",
		SSLEnabled:        false,
		MaxOpenConnection: 0,
		MaxIdleConnection: 0,
	})

	if err != nil {
		return xerrs.Wrap(err, "open database connection")
	}

	checkCtx, checkCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer checkCancel()

	if err := pqx.CheckConnection(checkCtx, 5, db); err != nil {
		return xerrs.Wrap(err, "checking database connection")
	}

	source := sqlize.NewSourceFromDir("vars/migrations")
	repository := sqlize.NewPostgreRepository(db, "sqlize_migrations")

	migrator := sqlize.NewMigrator(source, repository)

	switch args[0] {
	case "inspect":
	case "create":
		if len(args) < 2 {
			return xerrs.New("migration name is required")
		}

		_, err := migrator.New(args[1])
		if err != nil {
			return xerrs.Wrap(err, "exec create command")
		}
	case "status":
		if err := migrator.Status(context.Background()); err != nil {
			return xerrs.Wrap(err, "exec status command")
		}

		return nil
	case "init":
		if err := migrator.Init(context.Background()); err != nil {
			return xerrs.Wrap(err, "exec init command")
		}

		return nil
	case "up":
		if err := migrator.Up(context.Background()); err != nil {
			return xerrs.Wrap(err, "exec migrate command")
		}

		return nil
	case "undo":
		if err := migrator.Undo(context.Background()); err != nil {
			return xerrs.Wrap(err, "exec undo command")
		}

		return nil
	case "down":
		if err := migrator.Down(context.Background()); err != nil {
			return xerrs.Wrap(err, "exec undo command")
		}
		return nil
	case "help":
		fallthrough
	default:
		fmt.Print(usage)
	}

	return nil
}

const usage = `
sqlize tool help

sqlize <command> [<args>]

Command:
  help                  show this help.
  inspect               show configuration.

  up                    apply migration.
  init                  initialize migration table.
  undo                  undo last migration.
  down                  reset migration history into initial state.
  status                show migration status.
  create <name>         create a new migration file.
`
