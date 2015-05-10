package db

type Saver interface {
	Connect(location string) error
	Read(location string, key []byte) ([]byte, error)
	Save(location string, data []byte) (interface{}, error)
}



