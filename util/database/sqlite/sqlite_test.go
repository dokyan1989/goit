package sqlite

import (
	"database/sql"
	"errors"
	"os"
	"testing"
)

func TestCRUD(t *testing.T) {
	dsn := "dummy.db"

	db, err := New(dsn)
	if err != nil {
		t.Errorf("failed to open db, err = %v", err)
		return
	}
	defer func() {
		err := os.Remove(dsn)
		if err != nil {
			t.Fatal(err)
		}
	}()

	var dummySql = `
CREATE TABLE IF NOT EXISTS dummy (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(50) NOT NULL
);`

	_, err = db.Exec(dummySql)
	if err != nil {
		t.Errorf("failed to execute query, err = %v", err)
		return
	}

	err = insert(db)
	if err != nil {
		t.Errorf("failed to insert(), err = %v", err)
		return
	}

	dummies, err := get(db)
	if err != nil {
		t.Errorf("failed to get(), err = %v", err)
		return
	}

	if len(dummies) != 1 {
		t.Errorf("dummies length = %v, but want %v", len(dummies), 1)
	}

	err = update(db)
	if err != nil {
		t.Errorf("failed to update(), err = %v", err)
		return
	}

	err = delete(db)
	if err != nil {
		t.Errorf("failed to delete(), err = %v", err)
		return
	}
}

type dummy struct {
	ID   int64
	Name string
}

func get(db *sql.DB) ([]dummy, error) {
	var dummies []dummy

	rows, err := db.Query("SELECT * FROM dummy")
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

func insert(db *sql.DB) error {
	result, err := db.Exec("INSERT INTO dummy(name) VALUES (?)", "foo")
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("no row inserted")
	}

	return nil
}

func update(db *sql.DB) error {
	result, err := db.Exec("UPDATE dummy SET name = ? WHERE name = 'foo'", "bar")
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("no row updated")
	}

	return nil
}

func delete(db *sql.DB) error {
	result, err := db.Exec("DELETE FROM dummy WHERE name = 'bar'")
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("no row deleted")
	}

	return nil
}
