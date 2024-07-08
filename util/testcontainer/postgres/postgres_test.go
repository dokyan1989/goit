package postgres

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	ctx := context.Background()

	c, err := Run(ctx,
		WithDBName("foo"),
		WithDBUser("user"),
		WithDBPassword("password"),
	)
	if err != nil {
		t.Errorf("failed to Run() error = %v", err)
		return
	}

	dbpool, err := c.OpenDB(ctx)
	if err != nil {
		t.Errorf("failed to OpenDB() error = %v", err)
		return
	}

	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	sourceURL := fmt.Sprintf("file://%s", filepath.Join(workingDir, "./"))
	err = c.Migrate(ctx, sourceURL)
	if err != nil {
		t.Errorf("failed to Migrate() error = %v", err)
		return
	}

	var tableCount int
	err = dbpool.QueryRow(ctx, "SELECT count(*) FROM pg_tables WHERE schemaname = 'public'").Scan(&tableCount)
	if err != nil {
		t.Errorf("failed to dbpool.QueryRow() error = %v", err)
		return
	}

	if tableCount != 3 {
		t.Errorf("tableCount = %v but want %v", tableCount, 3)
		return
	}

	err = c.TruncateAllTables(ctx)
	if err != nil {
		t.Errorf("failed to TruncateAllTables() error = %v", err)
		return
	}

	err = c.DropAllTables(ctx)
	if err != nil {
		t.Errorf("failed to DropAllTables() error = %v", err)
		return
	}

	err = dbpool.QueryRow(ctx, "SELECT count(*) FROM pg_tables WHERE schemaname = 'public'").Scan(&tableCount)
	if err != nil {
		t.Errorf("failed to dbpool.QueryRow() error = %v", err)
		return
	}

	if tableCount != 0 {
		t.Errorf("tableCount = %v but want %v", tableCount, 0)
		return
	}

	err = c.Terminate(ctx)
	if err != nil {
		t.Errorf("failed to Terminate() error = %v", err)
		return
	}
}
