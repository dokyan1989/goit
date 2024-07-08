package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/dokyan1989/goit/migration"
	"github.com/dokyan1989/goit/misc/t/container/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	postgresC  *postgres.Container
	db         *pgxpool.Pool
	workingDir string
)

func runMain(ctx context.Context, m *testing.M) (code int, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c, err := postgres.Run(ctx,
		postgres.WithDBName("auth"),
		postgres.WithDBUser("user"),
		postgres.WithDBPassword("password"),
	)
	if err != nil {
		return 0, err
	}
	defer c.Terminate(ctx)
	postgresC = c // postgresC is defined as a package-level variable.

	d, err := c.OpenPGXDB(ctx)
	if err != nil {
		return 0, err
	}
	defer d.Close()
	db = d // db is defined as a package-level variable.

	wd, err := os.Getwd()
	if err != nil {
		return 0, fmt.Errorf("get working dir: %v", err)
	}
	workingDir = wd // workingDir is defined as a package-level variable.

	err = postgresC.MigrateFS(ctx, migration.AuthFS, migration.PathAuth)
	if err != nil {
		return 0, err
	}

	// m.Run() executes the regular, user-defined test functions.
	// Any defer statements that have been made will be run after m.Run()
	// completes.
	return m.Run(), nil
}

func TestMain(m *testing.M) {
	code, err := runMain(context.Background(), m)
	if err != nil {
		// Failure messages should be written to STDERR, which log.Fatal uses.
		log.Fatal(err)
	}

	// NOTE: defer statements do not run past here due to os.Exit terminating the process.
	os.Exit(code)
}
