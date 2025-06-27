package storage

const (
	STORAGE_TYPE_BADGER_DB = "badger"
	STORAGE_TYPE_MEMORY_DB = "memory"
)

type Storage interface {
	Get([]byte) ([]byte, error)
	Put([]byte, []byte) error
	Has([]byte) bool
	Delete([]byte) error
	BatchPut([][2][]byte) error
	Close() error
	Open() error
	GetIterator() Iterator
	GetSnapShot(string) SnapShot
}

func LoadDb(dbPath string, dbType string) (Storage, error) {
	var db Storage
	var err error
	switch dbType {
	case STORAGE_TYPE_BADGER_DB:
		db, err = NewBadgerDB(
			dbPath,
		)
		// case STORAGE_TYPE_MEMORY_DB:
		// 	db = NewMemoryDb()
	}
	return db, err
}
