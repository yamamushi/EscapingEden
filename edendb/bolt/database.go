package bolt

import (
	"errors"
	"github.com/asdine/storm/v3" // We use storm as a wrapper for boltdb, it's much easier to use than the native boltdb
	"github.com/yamamushi/EscapingEden/edendb"
	"sync"
)

type BoltDB struct {
	edendb.Database
	Path string // The path to the database file.

	queryMutex sync.Mutex
}

/*
	BoltDB can return the following errors (as known so far):

- already exists
- not found
- invalid argument

*/

func NewBoltDB(path string) (*BoltDB, error) {
	db := &BoltDB{Path: path}
	// Open [path] data file.
	// It will be created if it doesn't exist.
	// We're not really doing anything with the db here other than making sure the file exists.
	// If we hit an error here, nothing else is going to work either, so we can catch it and quit.
	boltDB, err := storm.Open(db.Path)
	if err != nil {
		return nil, err
	}
	defer boltDB.Close()

	db.Type = edendb.DatabaseTypeID_BoltDB
	return db, nil
}

func (db *BoltDB) AddRecord(collectionName string, value interface{}) error {
	db.queryMutex.Lock()
	defer db.queryMutex.Unlock()

	boltDB, err := storm.Open(db.Path)
	if err != nil {
		return err
	}
	defer boltDB.Close()

	collection := boltDB.From(collectionName)

	return collection.Save(value)
}

func (db *BoltDB) UpdateRecord(collectionName string, value interface{}) error {
	db.queryMutex.Lock()
	defer db.queryMutex.Unlock()

	boltDB, err := storm.Open(db.Path)
	if err != nil {
		return err
	}
	defer boltDB.Close()

	collection := boltDB.From(collectionName)

	return collection.Update(value)
}

func (db *BoltDB) UpdateField(collectionName string, field string, value interface{}, target interface{}) error {
	db.queryMutex.Lock()
	defer db.queryMutex.Unlock()

	boltDB, err := storm.Open(db.Path)
	if err != nil {
		return err
	}
	defer boltDB.Close()

	collection := boltDB.From(collectionName)

	return collection.UpdateField(target, field, value)
}

func (db *BoltDB) AddIfNotExists(collectionName string, value interface{}) error {
	db.queryMutex.Lock()
	defer db.queryMutex.Unlock()

	boltDB, err := storm.Open(db.Path)
	if err != nil {
		return err
	}
	defer boltDB.Close()

	return errors.New("method AddIfNotExists not implemented")
}

func (db *BoltDB) RemoveRecord(collectionName string, value interface{}) error {
	db.queryMutex.Lock()
	defer db.queryMutex.Unlock()

	boltDB, err := storm.Open(db.Path)
	if err != nil {
		return err
	}
	defer boltDB.Close()

	collection := boltDB.From(collectionName)
	return collection.DeleteStruct(value)
}

// TODO - This doesn't actually reduce the file size it seems, so we need to figure out a DB migration strategy
// That will let us copy the data to a new file so we can replace it manually and then delete the old file.

func (db *BoltDB) RemoveCollection(collectionName string) error {
	db.queryMutex.Lock()
	defer db.queryMutex.Unlock()

	boltDB, err := storm.Open(db.Path)
	if err != nil {
		return err
	}
	defer boltDB.Close()

	return boltDB.Drop(collectionName)
}

func (db *BoltDB) One(collectionName string, field string, value interface{}, output interface{}) error {
	db.queryMutex.Lock()
	defer db.queryMutex.Unlock()

	boltDB, err := storm.Open(db.Path)
	if err != nil {
		return err
	}
	defer boltDB.Close()

	collection := boltDB.From(collectionName)

	return collection.One(field, value, output)
}
