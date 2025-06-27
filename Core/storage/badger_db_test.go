package storage

import (
	"encoding/hex"
	fmt "fmt"
	"log/slog"
	"reflect"
	"testing"

	"FichainCore/common"
)

var bdb *BadgerDB

func init() {
	path := "./bdb_test"
	bdb, _ = NewBadgerDB(path)
}

func TestBadgerDB_Get(t *testing.T) {
	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Get",
			args: args{
				key: common.FromHex("0x0000000000000000000000000000000000000001"),
			},
			want:    common.FromHex("0x0000000000000000000000000000000000000002"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bdb.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("BadgerDB.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BadgerDB.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBadgerDB_Put(t *testing.T) {
	type args struct {
		key   []byte
		value []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Set",
			args: args{
				key:   common.FromHex("0x0000000000000000000000000000000000000001"),
				value: common.FromHex("0x0000000000000000000000000000000000000002"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := bdb.Put(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("BadgerDB.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBadgerDB_Has(t *testing.T) {
	type args struct {
		key []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Has",
			args: args{
				key: common.FromHex("0x0000000000000000000000000000000000000002"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bdb.Has(tt.args.key); got != tt.want {
				t.Errorf("BadgerDB.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBadgerDB_Delete(t *testing.T) {
	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Has",
			args: args{
				key: common.FromHex("0x0000000000000000000000000000000000000001"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := bdb.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("BadgerDB.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBadgerDB_BatchPut(t *testing.T) {
	type args struct {
		kvs [][2][]byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test batch put",
			args: args{
				kvs: [][2][]byte{
					{
						common.FromHex("0x0000000000000000000000000000000000000001"),
						common.FromHex("0x0000000000000000000000000000000000000002"),
					},
					{
						common.FromHex("0x0000000000000000000000000000000000000002"),
						common.FromHex("0x0000000000000000000000000000000000000003"),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := bdb.BatchPut(tt.args.kvs); (err != nil) != tt.wantErr {
				t.Errorf("BadgerDB.BatchPut() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBadgerDB_GetSnapShot(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "Test get snapshot and iter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapShot := bdb.GetSnapShot("./test_snap_shot")
			testV, _ := snapShot.(*BadgerDB).Get(
				common.FromHex("0x0000000000000000000000000000000000000002"),
			)
			slog.Info(fmt.Sprintf("value %v", hex.EncodeToString(testV)))
			iter := snapShot.GetIterator()
			slog.Info(fmt.Sprintf("iter %v", iter))
			for iter.Next() {
				slog.Info(
					fmt.Sprintf(
						"key %v & value %v",
						hex.EncodeToString(iter.Key()),
						hex.EncodeToString(iter.Value()),
					),
				)
			}
		})
	}
}
