package storage

type Iterator interface {
	Next() bool
	Key() []byte
	Value() []byte
	Release()
	Error() error
}
