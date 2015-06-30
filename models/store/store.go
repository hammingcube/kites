package store

import "errors"

var (
	ErrNotFound = errors.New("Not found")
)

type Store interface {
	Open(dbName string) error
	Close()
	Get(bucketName, key []byte) ([]byte, error)
	GetAll(bucketName []byte) ([][]byte, error)
	Post(bucketName, key, value []byte) error
	Delete(bucketName, key []byte) error
}
