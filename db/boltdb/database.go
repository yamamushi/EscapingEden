package boltdb

import (
	"github.com/yamamushi/EscapingEden/db"
	bolt "go.etcd.io/bbolt"
)

type BoltDB struct {
	db.Database
	Path string // The path to the database file.
}

func NewBoltDB(path string) (*BoltDB, error) {
	// Open [path] data file.
	// It will be created if it doesn't exist.
	boltdb, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	defer boltdb.Close()

	output := &BoltDB{Path: path}
	output.Type = db.DatabaseTypeID_BoltDB
	return output, nil
}

func (d *BoltDB) CreateCollection(collectionName string) error {
	return nil
}
