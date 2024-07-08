package postgres

import (
	"context"
	"errors"
	"testing"

	tc "github.com/dokyan1989/goit/util/testcontainer/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCRUD(t *testing.T) {
	ctx := context.Background()

	postgresC, err := tc.Run(ctx,
		tc.WithDBName("dummy"),
		tc.WithDBUser("foo"),
		tc.WithDBPassword("secret"),
	)
	if err != nil {
		t.Errorf("failed to run testcontainer, err = %v", err)
		return
	}
	defer postgresC.Terminate(ctx)

	db, err := New(ctx, postgresC.GetConnString())
	if err != nil {
		t.Errorf("failed to open db, err = %v", err)
		return
	}

	var dummySql = `
CREATE TABLE dummy (
  id SERIAL PRIMARY KEY,
  name varchar(50) NOT NULL
);`

	_, err = db.Exec(ctx, dummySql)
	if err != nil {
		t.Errorf("failed to execute query, err = %v", err)
		return
	}

	err = insert(ctx, db)
	if err != nil {
		t.Errorf("failed to insert(), err = %v", err)
		return
	}

	dummies, err := get(ctx, db)
	if err != nil {
		t.Errorf("failed to get(), err = %v", err)
		return
	}

	if len(dummies) != 1 {
		t.Errorf("dummies length = %v, but want %v", len(dummies), 1)
	}

	err = update(ctx, db)
	if err != nil {
		t.Errorf("failed to update(), err = %v", err)
		return
	}

	err = delete(ctx, db)
	if err != nil {
		t.Errorf("failed to delete(), err = %v", err)
		return
	}
}

type dummy struct {
	ID   int64
	Name string
}

func get(ctx context.Context, db *pgxpool.Pool) ([]dummy, error) {
	var dummies []dummy

	rows, err := db.Query(ctx, "SELECT * FROM dummy")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d dummy
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, err
		}
		dummies = append(dummies, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dummies, nil
}

func insert(ctx context.Context, db *pgxpool.Pool) error {
	// https://stackoverflow.com/a/74698904
	tag, err := db.Exec(ctx, "INSERT INTO dummy(name) VALUES ($1)", "foo")
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("no row inserted")
	}

	return nil
}

func update(ctx context.Context, db *pgxpool.Pool) error {
	tag, err := db.Exec(ctx, "UPDATE dummy SET name = $1 WHERE name = 'foo'", "bar")
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("no row updated")
	}

	return nil
}

func delete(ctx context.Context, db *pgxpool.Pool) error {
	tag, err := db.Exec(ctx, "DELETE FROM dummy WHERE name = 'bar'")
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("no row deleted")
	}

	return nil
}
