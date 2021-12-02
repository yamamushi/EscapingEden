package boltdb

import (
	"errors"
	"github.com/asdine/storm/v3" // We use storm as a wrapper for boltdb, it's much easier to use than the native boltdb
	"github.com/yamamushi/EscapingEden/db"
	"strings"
)

type BoltDB struct {
	db.Database
	Path string // The path to the database file.
	bolt storm.DB
}

func NewBoltDB(path string) (*BoltDB, error) {
	// Open [path] data file.
	// It will be created if it doesn't exist.
	boltDB, err := storm.Open(path)
	if err != nil {
		return nil, err
	}
	defer boltDB.Close()

	output := &BoltDB{Path: path}
	output.Type = db.DatabaseTypeID_BoltDB
	return output, nil
}

func (bd *BoltDB) GetNodeForPath(path string) (*storm.DB, storm.Node, error) {
	if path == "" {
		return nil, nil, errors.New("collection name cannot be empty")
	}

	boltDB, err := storm.Open(bd.Path)
	if err != nil {
		return nil, nil, err
	}
	defer boltDB.Close()

	// split collectionName into nested paths
	paths := strings.Split(path, "/")
	node := boltDB.From(paths[0])
	for i, path := range paths {
		if i > 0 {
			node = node.From(path)
		}
	}
	return boltDB, node, nil
}

func (bd *BoltDB) CreateCollection(collectionName string) error {

	return nil
}

func (bd *BoltDB) Update(input interface{}) error {
	boltDB, err := storm.Open(bd.Path)
	if err != nil {
		return err
	}
	defer boltDB.Close()

	return boltDB.Update(input)
}

func (bd *BoltDB) UpdateField(data interface{}, fieldName string, fieldValue interface{}) error {
	boltDB, err := storm.Open(bd.Path)
	if err != nil {
		return err
	}
	defer boltDB.Close()
	return boltDB.UpdateField(data, fieldName, fieldValue)
}

func (bd *BoltDB) Init(input interface{}) error {
	boltDB, err := storm.Open(bd.Path)
	if err != nil {
		return err
	}
	defer boltDB.Close()
	return boltDB.Init(input)
}

func (bd *BoltDB) UpdateNested(path string, input interface{}) error {
	boltDB, node, err := bd.GetNodeForPath(path)
	if err != nil {
		return err
	}
	defer boltDB.Close()

	return node.Update(input)
}

func (bd *BoltDB) UpdateFieldNested(path string, data interface{}, fieldName string, fieldValue interface{}) error {
	boltDB, node, err := bd.GetNodeForPath(path)
	if err != nil {
		return err
	}
	defer boltDB.Close()
	return node.UpdateField(data, fieldName, fieldValue)
}

func (bd *BoltDB) InitNested(path string, input interface{}) error {
	boltDB, node, err := bd.GetNodeForPath(path)
	if err != nil {
		return err
	}
	defer boltDB.Close()
	return node.Init(input)
}
