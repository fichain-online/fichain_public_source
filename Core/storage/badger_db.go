package storage

import (
	"encoding/hex"
	"errors"
	fmt "fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

type BadgerDB struct {
	db     *badger.DB
	closed bool
	path   string
	mu     *sync.RWMutex
}

// badger iterator
type BadgerDBIterator struct {
	db          *badger.DB
	txn         *badger.Txn
	iterator    *badger.Iterator
	currentItem *badger.Item
}

func NewBadgerDB(path string) (*BadgerDB, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}

	// Run garbage collection every hour.
	go func() {
		for range time.Tick(time.Minute) {
			db.RunValueLogGC(0.7)
		}
	}()

	return &BadgerDB{
		db:     db,
		closed: false,
		path:   path,
		mu:     &sync.RWMutex{},
	}, nil
}

func NewBadgerDBWithSnapshot(path string, snapshotFilePath string) (*BadgerDB, error) {
	// Open DB
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	// Load backup from snapshot file
	file, err := os.Open(snapshotFilePath)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to open snapshot file: %w", err)
	}
	defer file.Close()

	if err := db.Load(file, 10); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to load snapshot: %w", err)
	}

	// Run garbage collection every hour
	go func() {
		for range time.Tick(time.Minute) {
			db.RunValueLogGC(0.7)
		}
	}()

	return &BadgerDB{
		db:     db,
		closed: false,
		path:   path,
		mu:     &sync.RWMutex{},
	}, nil
}

func (bdb *BadgerDB) Get(key []byte) ([]byte, error) {
	if bdb.closed {
		return nil, errors.New("database is closed")
	}

	var value []byte
	err := bdb.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})
	})
	if err != nil && err != badger.ErrKeyNotFound {
		return nil, err
	}

	return value, nil
}

func (bdb *BadgerDB) Put(key, value []byte) error {
	bdb.mu.Lock()
	defer bdb.mu.Unlock()

	if bdb.closed {
		return errors.New("database is closed")
	}

	err := bdb.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})

	return err
}

func (bdb *BadgerDB) Has(key []byte) bool {
	if bdb.closed {
		return false
	}

	err := bdb.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err != nil {
			return err
		}
		return nil
	})

	return err == nil
}

func (bdb *BadgerDB) Delete(key []byte) error {
	bdb.mu.Lock()
	defer bdb.mu.Unlock()

	if bdb.closed {
		return errors.New("database is closed")
	}

	err := bdb.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})

	return err
}

func (bdb *BadgerDB) BatchPut(kvs [][2][]byte) error {
	bdb.mu.Lock()
	defer bdb.mu.Unlock()

	if bdb.closed {
		return errors.New("database is closed")
	}

	err := bdb.db.Update(func(txn *badger.Txn) error {
		for _, kv := range kvs {
			if err := txn.Set(kv[0], kv[1]); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (bdb *BadgerDB) Open() error {
	if !bdb.closed {
		return nil
	}
	var err error
	bdb.db, err = badger.Open(badger.DefaultOptions(bdb.path))
	if err != nil {
		return err
	}
	bdb.closed = false
	return nil
}

func (bdb *BadgerDB) Close() error {
	if bdb.closed {
		return nil
	}

	err := bdb.db.Close()
	bdb.closed = true

	return err
}

func (bdb *BadgerDB) GetSnapShot(
	snapShotDir string,
) SnapShot {
	bdb.mu.Lock()
	defer bdb.mu.Unlock()

	if bdb.closed {
		slog.Error("attempt to snapshot closed DB")
		return nil
	}

	// Ensure directory exists
	if err := os.MkdirAll(snapShotDir, 0755); err != nil {
		slog.Error("failed to create snapshot directory", slog.Any("error", err))
		return nil
	}

	// Generate a filename like: snapshot-20250509T150405.bak
	now := time.Now()
	filename := fmt.Sprintf("snapshot-%s.bak", now.Format("20060102T150405"))
	snapShotPath := filepath.Join(snapShotDir, filename)

	// Create snapshot file inside the directory
	file, err := os.Create(snapShotPath)
	if err != nil {
		slog.Error("failed to create snapshot file", slog.Any("error", err))
		return nil
	}
	defer file.Close()

	_, err = bdb.db.Backup(file, 0)
	if err != nil {
		slog.Error("failed to write snapshot", slog.Any("error", err))
		return nil
	}

	snapShot, err := NewBadgerDBWithSnapshot(snapShotDir, snapShotPath)
	if err != nil {
		slog.Error("failed to read snapshot", slog.Any("error", err))
		return nil
	}
	return snapShot
}

func (bdb *BadgerDB) GetIterator() Iterator {
	return NewBadgerDBIterator(bdb)
}

func (bdb *BadgerDB) Release() {
	bdb.db.DropAll()
	bdb.db.Close()
}

func NewBadgerDBIterator(bdb *BadgerDB) *BadgerDBIterator {
	txn := bdb.db.NewTransaction(true)
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 1000
	iterator := txn.NewIterator(opts)
	iterator.Rewind()
	return &BadgerDBIterator{
		db:       bdb.db,
		txn:      txn,
		iterator: iterator,
	}
}

func (bIter *BadgerDBIterator) Next() bool {
	if bIter.iterator.Valid() {
		item := bIter.iterator.Item()
		k := item.Key()
		item.Value(func(v []byte) error {
			fmt.Printf("key=%s, value=%s\n", hex.EncodeToString(k), hex.EncodeToString(v))
			return nil
		})
		bIter.currentItem = bIter.iterator.Item()
		bIter.iterator.Next()
		return true
	}
	// Rewind the iterator to the beginning of the database.
	return false
}

func (bIter *BadgerDBIterator) Key() []byte {
	return bIter.currentItem.Key()
}

func (bIter *BadgerDBIterator) Value() []byte {
	v, _ := bIter.currentItem.ValueCopy(nil)
	return v
}

func (bIter *BadgerDBIterator) Release() {
}

func (bIter *BadgerDBIterator) Error() error {
	return nil
}
