package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const DB_HOST = "127.0.0.1"

var dbSession *mgo.Session

type Users struct {
	db string
	coll string 
}

type User struct {
	Name	string
	Email	string
}

func (users *Users) Save(u *User) error {
	var err error
	session := dbSession.Copy()
    defer session.Close()
    if c := session.DB(users.db).C(users.coll); c != nil {
    	err = c.Insert(u)
    }
    return err
}

func (users *Users) Read(email string) (*User, error) {
	var err error
	session := dbSession.Copy()
    defer session.Close()
    if c := session.DB(users.db).C(users.coll); c != nil {
    	result := &User{}
    	err = c.Find(bson.M{"email": email}).One(&result)
    	return result, err
    }
    return nil, err
}

func main() {
	var err error
	dbSession, err = mgo.Dial(DB_HOST)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	dbSession.SetMode(mgo.Monotonic, true)
	defer dbSession.Close()

	users := &Users{"myfavdb", "users"}
	users.Save(&User{"John Doe", "jdoe@example.com"})
	email := "jdoe@example.com"
	user, err := users.Read(email)
	if err != nil {
		log.Printf("Error when finding user %s: %v", email, err)
	} else {
		log.Printf("%s", user)
	}

}