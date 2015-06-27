package users

import "encoding/json"
import "github.com/maddyonline/kites/models/store"

type User struct {
	Id       	string
	Username 	string
	Email    	string
	Meta		*Data
}

type Data struct {
	GistsIds	[]string
}


type Store interface {
	Open(dbName string) error
	Close()
	Get(bucketName, key string) (*User, error)
	GetAll(bucketName string) (*[]User, error)
	Post(bucketName string, key string, u *User) error
	Delete(bucketName string, key string) error
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

func (i *Impl) Get(bucketName, key string) (*User, error) {
	bytes, err := i.db.Get([]byte(bucketName), []byte(key))
	if err != nil {
		return nil, err
	}
	g, err := FromBytes(bytes)
	return g, err
}

func (i *Impl) Post(bucketName, key string, g *User) error {
	data, err := ToBytes(g)
	if err != nil {
		return err
	}
	err = i.db.Post([]byte(bucketName), []byte(key), data)
	return err
}

func (i *Impl) Delete(bucketName, key string) error {
	err := i.db.Delete([]byte(bucketName), []byte(key))
	return err
}

func (i *Impl) GetAll(bucketName string) (*[]User, error) {
	data, err := i.db.GetAll([]byte(bucketName))
	if err != nil {
		return nil, err
	}
	users := []User{}
	for _, bytes := range data {
		g, err := FromBytes(bytes)
		if err != nil {
			return nil, err
		}
		users = append(users, *g)
	}
	return &users, nil
}

func (i *Impl) Open(dbName string) error {
	return i.db.Open(dbName)
}

func (i *Impl) Close() {
	i.db.Close()
}
