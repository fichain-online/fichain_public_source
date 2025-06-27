package database

import (
	"path/filepath"
	"testing"
)

func createTestDB(t *testing.T) *BadgerDB {
	t.Helper()
	dir := t.TempDir()
	db, err := NewBadgerDB(filepath.Join(dir, "badger"))
	if err != nil {
		t.Fatalf("failed to create test DB: %v", err)
	}
	return db
}

func TestBadgerDB_PutGetHasDelete(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	key := []byte("my-key")
	value := []byte("my-value")

	// Test Put
	err := db.Put(key, value)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Test Get
	got, err := db.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(got) != string(value) {
		t.Errorf("expected %s, got %s", value, got)
	}

	// Test Has
	has, err := db.Has(key)
	if err != nil {
		t.Fatalf("Has failed: %v", err)
	}
	if !has {
		t.Errorf("expected key to exist")
	}

	// Test Delete
	err = db.Delete(key)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Test Has after delete
	has, err = db.Has(key)
	if err != nil {
		t.Fatalf("Has after delete failed: %v", err)
	}
	if has {
		t.Errorf("expected key to be deleted")
	}
}

func TestBadgerBatch_WriteAndReset(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	batch := db.NewBatch()

	key1 := []byte("batch-key-1")
	val1 := []byte("value1")
	key2 := []byte("batch-key-2")
	val2 := []byte("value2")

	// Add to batch
	if err := batch.Put(key1, val1); err != nil {
		t.Fatalf("batch Put key1 failed: %v", err)
	}
	if err := batch.Put(key2, val2); err != nil {
		t.Fatalf("batch Put key2 failed: %v", err)
	}

	// Ensure ValueSize is correct
	expectedSize := len(key1) + len(val1) + len(key2) + len(val2)
	if batch.ValueSize() != expectedSize {
		t.Errorf("expected batch size %d, got %d", expectedSize, batch.ValueSize())
	}

	// Write batch
	if err := batch.Write(); err != nil {
		t.Fatalf("batch Write failed: %v", err)
	}

	// Confirm values in DB
	got1, err := db.Get(key1)
	if err != nil || string(got1) != string(val1) {
		t.Errorf("expected value1, got %s (err=%v)", got1, err)
	}
	got2, err := db.Get(key2)
	if err != nil || string(got2) != string(val2) {
		t.Errorf("expected value2, got %s (err=%v)", got2, err)
	}

	// Reset batch and reuse
	batch.Reset()
	if batch.ValueSize() != 0 {
		t.Errorf("expected reset batch size 0, got %d", batch.ValueSize())
	}

	key3 := []byte("batch-key-3")
	val3 := []byte("value3")
	if err := batch.Put(key3, val3); err != nil {
		t.Fatalf("batch Put key3 failed after reset: %v", err)
	}
	if err := batch.Write(); err != nil {
		t.Fatalf("batch Write after reset failed: %v", err)
	}

	got3, err := db.Get(key3)
	if err != nil || string(got3) != string(val3) {
		t.Errorf("expected value3, got %s (err=%v)", got3, err)
	}
}
