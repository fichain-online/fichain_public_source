package storage

import (
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"sync"
)

type MemoryDB struct {
	db map[string][]byte
	sync.RWMutex
}

type MemoryDbIterator struct {
	memoryDB *MemoryDB
	keys     []string
	idx      int
	sync.RWMutex
}

func NewMemoryDbIterator(memoryDB *MemoryDB) *MemoryDbIterator {
	var keys []string
	memoryDB.RLock()
	for key := range memoryDB.db {
		keys = append(keys, key)
	}
	memoryDB.RUnlock()
	sort.Strings(keys)
	return &MemoryDbIterator{
		memoryDB: memoryDB,
		keys:     keys,
		idx:      0,
	}
}

func (mdb *MemoryDB) GetIterator() Iterator {
	mdb.RLock()
	defer mdb.RUnlock()
	return NewMemoryDbIterator(mdb)
}

func (mdb *MemoryDbIterator) Next() bool {
	mdb.Lock()
	defer mdb.Unlock()
	return mdb.idx < len(mdb.keys) && (func() bool { mdb.idx++; return true })()
}

func (mdb *MemoryDbIterator) Key() []byte {
	mdb.RLock()
	defer mdb.RUnlock()
	return []byte(mdb.keys[mdb.idx-1])
}

func (mdb *MemoryDbIterator) Value() []byte {
	mdb.RLock()
	mdb.memoryDB.RLock()
	defer mdb.RUnlock()
	defer mdb.memoryDB.RUnlock()
	return mdb.memoryDB.db[mdb.keys[mdb.idx-1]]
}

func (mdb *MemoryDbIterator) Release() {}

func (mdb *MemoryDbIterator) Error() error {
	return nil
}

func NewMemoryDb() *MemoryDB {
	return &MemoryDB{
		db: make(map[string][]byte),
	}
}

func (kv *MemoryDB) Get(key []byte) ([]byte, error) {
	kv.RLock()
	defer kv.RUnlock()
	strKey := string(key)
	if v, ok := kv.db[strKey]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("[MemKV] key not found: %s", hex.EncodeToString(key)))
}

func (kv *MemoryDB) Put(key, value []byte) error {
	kv.Lock()
	defer kv.Unlock()
	kv.db[string(key)] = value
	return nil
}

func (kv *MemoryDB) Has(key []byte) bool {
	kv.RLock()
	defer kv.RUnlock()
	_, ok := kv.db[string(key)]
	return ok
}

func (kv *MemoryDB) Delete(key []byte) error {
	kv.Lock()
	defer kv.Unlock()
	strKey := string(key)
	if _, ok := kv.db[strKey]; ok {
		delete(kv.db, strKey)
		return nil
	}
	return errors.New(fmt.Sprintf("[MemKV] key not found: %s", hex.EncodeToString(key)))
}

func (kv *MemoryDB) BatchPut(kvs [][2][]byte) error {
	kv.Lock()
	defer kv.Unlock()
	for _, kvp := range kvs {
		kv.db[string(kvp[0])] = kvp[1]
	}
	return nil
}

func (kv *MemoryDB) Close() error   { return nil }
func (kv *MemoryDB) Open() error    { return nil }
func (kv *MemoryDB) Compact() error { return nil }
func (kv *MemoryDB) Size() int      { return len(kv.db) }

func (kv *MemoryDB) GetSnapShot(path string) SnapShot {
	newMDB := NewMemoryDb()
	iter := kv.GetIterator()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		cKey := make([]byte, len(k))
		cValue := make([]byte, len(v))
		copy(cKey, k)
		copy(cValue, v)
		newMDB.Put(cKey, cValue)
	}
	return newMDB
}

func (kv *MemoryDB) Release() {}
