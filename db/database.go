package db

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

	CreateCollection(collectionName string) error

	Update(interface{}) error
	UpdateField(interface{}, string, interface{}) error
	Init(interface{}) error

	UpdateNested(string, interface{}) error
	UpdateFieldNested(string, interface{}, string, interface{}) error
	InitNested(string, interface{}) error
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

func (db *Database) CreateCollection(collectionName string) error {
	return nil
}

func (db *Database) Update(interface{}) error {
	return nil
}

func (db *Database) UpdateField(interface{}, string, interface{}) error {
	return nil
}

func (db *Database) Init(interface{}) error {
	return nil
}

func (db *Database) UpdateNested(string, interface{}) error {
	return nil
}

func (db *Database) UpdateFieldNested(string, interface{}, string, interface{}) error {
	return nil
}

func (db *Database) InitNested(string, interface{}) error {
	return nil
}
