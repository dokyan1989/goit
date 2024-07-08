package seeder

import (
	"context"
	"testing"

	"github.com/dokyan1989/goit/misc/t/container/postgres"
)

func MustRun(ctx context.Context, t *testing.T, c *postgres.Container, seedURL string) {
	t.Helper()

	if c == nil {
		t.Fatal("container is nil")
	}

	if err := c.Seed(ctx, seedURL); err != nil {
		t.Fatalf("seeding data: %v", err)
	}

	t.Cleanup(func() {
		if err := c.TruncateAllTables(ctx); err != nil {
			t.Fatalf("truncate all tables: %v", err)
		}
	})
}
