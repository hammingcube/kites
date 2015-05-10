package db

import (
	"gopkg.in/mgo.v2"
)

type Mongo struct {
	cfg	Configuration
	session	*mgo.Session
}

//type interface {
	//Connect() err
	//Save(location string, data []byte) (err, interface{})
	//Read(location string, key []byte) (err, []byte)
//}

type Configuration struct {
	host	string
}

func (m *Mongo) Connect() error {
	session, err := mgo.Dial(m.cfg.host)
	if err == nil {
		m.session = session
	}
	return err
}