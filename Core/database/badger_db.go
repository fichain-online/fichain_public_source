package database

import (
	"sync"

	badger "github.com/dgraph-io/badger/v4"
)

// BadgerDB implements the Database interface using Badger.
type BadgerDB struct {
	Path string
	db   *badger.DB
}

// NewBadgerDB initializes a new BadgerDB instance.
func NewBadgerDB(path string) (*BadgerDB, error) {
	opts := badger.DefaultOptions(path).WithLogger(nil) // disable logging
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerDB{
		db:   db,
		Path: path,
	}, nil
}

func (b *BadgerDB) Put(key []byte, value []byte) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (b *BadgerDB) Get(key []byte) ([]byte, error) {
	var valCopy []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		valCopy, err = item.ValueCopy(nil)
		return err
	})
	return valCopy, err
}

func (b *BadgerDB) Has(key []byte) (bool, error) {
	err := b.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		return err
	})
	if err == badger.ErrKeyNotFound {
		return false, nil
	}
	return err == nil, err
}

func (b *BadgerDB) Delete(key []byte) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (b *BadgerDB) IterateKeys(fn func(key, value []byte) error) error {
	return b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.KeyCopy(nil)
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			if err := fn(key, val); err != nil {
				return err
			}
		}
		return nil
	})
}

func (b *BadgerDB) Close() {
	b.db.Close()
}

func (b *BadgerDB) NewBatch() Batch {
	return &BadgerBatch{
		db:    b.db,
		wb:    b.db.NewWriteBatch(),
		size:  0,
		mutex: &sync.Mutex{},
	}
}

// BadgerBatch implements the Batch interface using Badger WriteBatch.
type BadgerBatch struct {
	db    *badger.DB
	wb    *badger.WriteBatch
	size  int
	mutex *sync.Mutex
}

func (b *BadgerBatch) Put(key []byte, value []byte) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	err := b.wb.Set(key, value)
	if err != nil {
		return err
	}
	b.size += len(key) + len(value)
	return nil
}

func (b *BadgerBatch) ValueSize() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.size
}

func (b *BadgerBatch) Write() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	err := b.wb.Flush()
	b.size = 0
	return err
}

func (b *BadgerBatch) Reset() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.wb.Cancel() // discard current
	b.wb = b.db.NewWriteBatch()
	b.size = 0
}
