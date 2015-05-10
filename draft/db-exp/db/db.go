package db

import "gopkg.in/mgo.v2"
import "time"

const DIAL_TIMEOUT = time.Duration(2)*time.Second

type Config struct {
	Host	string
}

type Database struct {
	Cfg	Config
	Session *mgo.Session
}

func NewConnection(host string) (*Database, error) {
	db := &Database{Cfg: Config{Host: host}}
	session, err := mgo.DialWithTimeout(host, DIAL_TIMEOUT)
	if err == nil {
		db.Session = session
	}
	return db, err
}
