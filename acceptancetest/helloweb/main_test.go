package helloweb

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/dokyan1989/goit/migration"
	"github.com/dokyan1989/goit/misc/t/container/postgres"
)

var (
	postgresC  *postgres.Container
	workingDir string
)

func runMain(ctx context.Context, m *testing.M) (code int, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c, err := postgres.Run(ctx,
		postgres.WithDBName("todo"),
		postgres.WithDBUser("user"),
		postgres.WithDBPassword("password"),
	)
	if err != nil {
		return 0, err
	}
	defer c.Terminate(ctx)
	postgresC = c // postgresC is defined as a package-level variable.

	os.Setenv("HELLO_DB_URL", postgresC.GetConnString())
	os.Setenv("HELLO_ENV", "local")

	wd, err := os.Getwd()
	if err != nil {
		return 0, fmt.Errorf("get working dir: %v", err)
	}
	workingDir = wd

	// err = postgresC.Migrate(ctx, fmt.Sprintf("file://%s", filepath.Join(workingDir, "../../internal/todo/migration")))
	err = postgresC.MigrateFS(ctx, migration.TodoFS, migration.PathTodo)
	if err != nil {
		return 0, err
	}

	/**
	|-------------------------------------------------------------------------
	| Launch test program
	|-----------------------------------------------------------------------*/
	cleanup, _, err := LaunchTestProgram("3000")
	if err != nil {
		return 0, err
	}
	defer cleanup()

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
