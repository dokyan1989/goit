package badger

import (
	"bytes"
	"encoding/gob"
	"os"
	"strconv"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/dokyan1989/goit/misc/snowflake"
)

func TestCRUDInMem(t *testing.T) {
	db, err := NewInMem()
	if err != nil {
		t.Errorf("failed to open db, err = %v", err)
		return
	}

	// --- test insert ---
	err = insert(db)
	if err != nil {
		t.Errorf("failed to insert(), err = %v", err)
		return
	}

	// --- test get ---
	dummies, err := get(db)
	if err != nil {
		t.Errorf("failed to get(), err = %v", err)
		return
	}

	if len(dummies) != 1 {
		t.Errorf("dummies length = %v, but want %v", len(dummies), 1)
		return
	}

	// --- test update ---
	err = update(db)
	if err != nil {
		t.Errorf("failed to update(), err = %v", err)
		return
	}

	dummies, err = get(db)
	if err != nil {
		t.Errorf("failed to get(), err = %v", err)
		return
	}

	if len(dummies) != 1 {
		t.Errorf("dummies length = %v, but want %v", len(dummies), 1)
		return
	}

	if dummies[0].Name != "bar" {
		t.Errorf("dummies[0].Name = %v, but want %v", dummies[0].Name, 1)
		return
	}

	// --- test delete ---
	err = delete(db)
	if err != nil {
		t.Errorf("failed to delete(), err = %v", err)
		return
	}

	dummies, err = get(db)
	if err != nil {
		t.Errorf("failed to get(), err = %v", err)
		return
	}

	if len(dummies) != 0 {
		t.Errorf("dummies length = %v, but want %v", len(dummies), 0)
		return
	}
}

func TestCRUDWithPath(t *testing.T) {
	path := "./sample"
	db, err := NewWithPath(path)
	if err != nil {
		t.Errorf("failed to open db, err = %v", err)
		return
	}
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// --- test insert ---
	err = insert(db)
	if err != nil {
		t.Errorf("failed to insert(), err = %v", err)
		return
	}

	// --- test get ---
	dummies, err := get(db)
	if err != nil {
		t.Errorf("failed to get(), err = %v", err)
		return
	}

	if len(dummies) != 1 {
		t.Errorf("dummies length = %v, but want %v", len(dummies), 1)
		return
	}

	// --- test update ---
	err = update(db)
	if err != nil {
		t.Errorf("failed to update(), err = %v", err)
		return
	}

	dummies, err = get(db)
	if err != nil {
		t.Errorf("failed to get(), err = %v", err)
		return
	}

	if len(dummies) != 1 {
		t.Errorf("dummies length = %v, but want %v", len(dummies), 1)
		return
	}

	if dummies[0].Name != "bar" {
		t.Errorf("dummies[0].Name = %v, but want %v", dummies[0].Name, 1)
		return
	}

	// --- test delete ---
	err = delete(db)
	if err != nil {
		t.Errorf("failed to delete(), err = %v", err)
		return
	}

	dummies, err = get(db)
	if err != nil {
		t.Errorf("failed to get(), err = %v", err)
		return
	}

	if len(dummies) != 0 {
		t.Errorf("dummies length = %v, but want %v", len(dummies), 0)
		return
	}
}

type dummy struct {
	ID   int64
	Name string
}

func insert(db *badger.DB) error {
	err := db.Update(func(txn *badger.Txn) error {
		// generate a snowflake id as key
		id, err := snowflake.NextID()
		if err != nil {
			return err
		}

		// create a dummy record "foo"
		d := dummy{ID: id, Name: "foo"}

		// encode the record to []byte
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err = enc.Encode(d)
		if err != nil {
			return err
		}

		// add the record to db
		err = txn.Set([]byte(strconv.FormatInt(id, 10)), buf.Bytes())
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func get(db *badger.DB) ([]dummy, error) {
	var dummies []dummy

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			id, err := strconv.ParseInt(string(item.Key()), 10, 64)
			if err != nil {
				return err
			}

			err = item.Value(func(v []byte) error {
				var d dummy
				dec := gob.NewDecoder(bytes.NewBuffer(v))
				err := dec.Decode(&d)
				if err != nil {
					return err
				}

				dummies = append(dummies, dummy{
					ID:   id,
					Name: d.Name,
				})
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return dummies, nil
}

func update(db *badger.DB) error {
	err := db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()

			id, err := strconv.ParseInt(string(k), 10, 64)
			if err != nil {
				return err
			}

			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			err = enc.Encode(dummy{ID: id, Name: "bar"})
			if err != nil {
				return err
			}

			err = txn.Set(k, buf.Bytes())
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func delete(db *badger.DB) error {
	err := db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()

			err := txn.Delete(k)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
