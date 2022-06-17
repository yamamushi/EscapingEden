package edendb

import (
	"io"
)

type DatabaseTypeID int

const (
	DatabaseTypeID_Unknown DatabaseTypeID = iota
	DatabaseTypeID_BoltDB
	DatabaseTypeID_MongoDB
)

func (dt *DatabaseTypeID) String() string {
	switch *dt {
	case DatabaseTypeID_Unknown:
		return "Unknown"
	case DatabaseTypeID_BoltDB:
		return "BoltDB"
	case DatabaseTypeID_MongoDB:
		return "MongoDB"
	default:
		return "Unknown"
	}
}

type DatabaseType interface {
	GetTypeID() DatabaseTypeID
	GetTypeName() string
	Init() error

	AddRecord(string, interface{}) error                        // Add record to collection, return error if failed
	UpdateRecord(string, interface{}) error                     // Update record in collection if it exists, otherwise add it
	UpdateField(string, string, interface{}, interface{}) error // Update field in record if it exists, otherwise add it
	AddIfNotExists(string, interface{}) error                   // Add record to collection if it doesn't exist

	RemoveRecord(string, interface{}) error // Remove a specific record from the given collection
	RemoveCollection(string) error          // Remove the given collection and all its records

	// One single search result
	// Collection, Field, Value, Output
	One(string, string, interface{}, interface{}) error

	// FindAll records in collection that match the interface.
	FindAll(string, []interface{}) error

	// FindAllByField finds all interfaces by - collection, field, search value
	FindAllByField(string, string, string, []interface{}) error

	// DumpDatabase dumps the entire database as []interface{}
	DumpDatabase([]interface{}) error

	// DumpCollection dumps the entire collection as []interface{}
	DumpCollection(string, []interface{}) error

	// DumpToWriter dumps the entire database to the writer
	DumpToWriter(io.Writer) error

	// DumpCollectionToWriter dumps the entire collection to the writer
	DumpCollectionToWriter(string, io.Writer) error
}
type Database struct {
	Type DatabaseTypeID
}

func (db *Database) GetTypeID() DatabaseTypeID {
	return db.Type
}

func (db *Database) GetTypeName() string {
	return db.Type.String()
}

func (db *Database) Init() error {
	return nil
}

func (db *Database) AddRecord(collectionName string, value interface{}) error {
	return nil
}

func (db *Database) UpdateRecord(collectionName string, value interface{}) error {
	return nil
}

func (db *Database) AddIfNotExists(collectionName string, value interface{}) error {
	return nil
}

func (db *Database) RemoveRecord(collectionName string, value interface{}) error {
	return nil
}

func (db *Database) RemoveCollection(collectionName string) error {
	return nil
}

func (db *Database) One(collectionName string, field string, value interface{}, output interface{}) error {
	return nil
}

func (db *Database) FindAll(collectionName string, output []interface{}) error {
	return nil
}

func (db *Database) FindAllByField(collectionName string, fieldName string, searchValue string, output []interface{}) error {
	return nil
}

func (db *Database) DumpDatabase(output []interface{}) error {
	return nil
}

func (db *Database) DumpCollection(collectionName string, output []interface{}) error {
	return nil
}

func (db *Database) DumpToWriter(writer io.Writer) error {
	return nil
}

func (db *Database) DumpCollectionToWriter(collectionName string, writer io.Writer) error {
	return nil
}
