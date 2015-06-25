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

func ToBytes()

type Store interface {
	Open(dbName, bucketName string) error
	Close()
	Get(key string) ([]byte, error)
	GetAll() ([][]byte, error)
	Post([]byte) error
	Delete(key string) error
}

type Impl struct {
	db store.Store 
}

func doIt() {
	db := &boltdb.Store{}
	gstore := &gists.Impl{db: db}
	gstore.Open("aa", "bb")
	defer gstore.Close()
}

func doIt2() {
	db := &mgo.Store{}
	gstore := &gists.Impl{db: db}
	gstore.Open("aa", "bb")
	defer gstore.Close()
}

func (i *Impl) Get(key string) (*Gist, error) {
		bytes, err := i.db.Get([]byte(key))
		if err != nil {
			return nil, err
		}
		g, err := FromBytes(bytes)
		return g, err
	}
}

func (i *Impl) Post(g *Gist) error {
		bytes, err := ToBytes(g)
		if err != nil {
			return err
		}
		err = i.db.Post(bytes)
		return err
	}
}

func (i *Impl) Delete(key string) error {
	err := i.db.Delete([]byte(key))
	return err
}

func (i *Impl) GetAll() (*[]Gist, error) {
		data, err := i.db.GetAll()
		if err != nil {
			return nil, err
		}
		gists := []Gist{}
		for _, bytes := range data {
			g, err := FromBytes(bytes)
			if err != nil {
				return nil, error
			}
			gists = append(gists, g)
		}
		return &gists, nil
}

func (i *Impl) Open(dbName, bucketName string) error {
	return i.db.Open(dbName, bucketName)
}

func (i *Impl) Close() {
	i.db.Close()
}

type Store interface {
	Open(dbName, bucketName string) error
	Close()
	Get(key string) (*Gist, error)
	GetAll() ([]Gist, error)
	Post(u *Gist) error
	Delete(key string) error
}
