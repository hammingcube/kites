package gists

import "github.com/maddyonline/kites/draft/db-exp/bolt_db/users"

type Gist struct {
	Id    string
	Files map[string]File
	*users.User
}

type File struct {
	Content  []byte
	Name     string
	Language string
}

func FromBytes([]byte) *Gist {
	return &Gist{}
}


type Store interface {
	Open(dbName, bucketName string) error
	Close()
	Get(key string) (*Gist, error)
	GetAll() ([]Gist, error)
	Post(u *Gist) error
	Delete(key string) error
}
