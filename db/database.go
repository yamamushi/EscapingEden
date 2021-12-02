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
