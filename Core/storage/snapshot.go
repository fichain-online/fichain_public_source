package storage

type SnapShot interface {
	GetIterator() Iterator
	Release()
}
