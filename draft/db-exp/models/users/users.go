package users

import "encoding/json"
import "github.com/maddyonline/kites/draft/db-exp/models/store"

type User struct {
	Id       string
	Username string
	Email    string
}

type Store interface {
	Open(dbName, bucketName string) error
	Close()
	Get(key string) (*User, error)
	GetAll() (*[]User, error)
	Post(u *User) error
	Delete(key string) error
}

func FromBytes(data []byte) (*User, error) {
	g := &User{}
	err := json.Unmarshal(data, g)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func ToBytes(g *User) ([]byte, error) {
	return json.Marshal(g)
}

type Impl struct {
	db store.Store
}

func NewStore(db store.Store) Store {
	i := &Impl{db: db}
	return i
}

func (i *Impl) Get(key string) (*User, error) {
	bytes, err := i.db.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	g, err := FromBytes(bytes)
	return g, err
}

func (i *Impl) Post(g *User) error {
	key := []byte(g.Id)
	data, err := ToBytes(g)
	if err != nil {
		return err
	}
	err = i.db.Post(key, data)
	return err
}

func (i *Impl) Delete(key string) error {
	err := i.db.Delete([]byte(key))
	return err
}

func (i *Impl) GetAll() (*[]User, error) {
	data, err := i.db.GetAll()
	if err != nil {
		return nil, err
	}
	gists := []User{}
	for _, bytes := range data {
		g, err := FromBytes(bytes)
		if err != nil {
			return nil, err
		}
		gists = append(gists, *g)
	}
	return &gists, nil
}

func (i *Impl) Open(dbName, bucketName string) error {
	return i.db.Open(dbName, bucketName)
}

func (i *Impl) Close() {
	i.db.Close()
}
