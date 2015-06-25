package store

type Store interface {
	Open(dbName, bucketName string) error
	Close()
	Get([]byte) ([]byte, error)
	GetAll() ([][]byte, error)
	Post([]byte, []byte) error
	Delete([]byte) error
}
