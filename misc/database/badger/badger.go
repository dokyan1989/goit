package badger

import (
	"github.com/dgraph-io/badger/v4"
)

func NewInMem() (*badger.DB, error) {
	opt := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewWithPath(path string) (*badger.DB, error) {
	opt := badger.DefaultOptions(path)
	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	return db, nil
}
