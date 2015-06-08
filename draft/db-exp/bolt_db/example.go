package main

import (
	"github.com/maddyonline/kites/draft/db-exp/bolt_db/users"
	"log"
)

func main() {
	store := &users.BoltDBStore{}
	DoIt(store)
}

func DoIt(store users.Store) {
	if err := store.Open("usersdb", "users"); err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	u1 := &users.User{"maddy", "maddy", "maddy@gmail.com"}
	u2 := &users.User{"mj", "mj", "mj@gmail.com"}
	err := store.Post(u1)
	err = store.Post(u2)
	if err != nil {
		log.Fatal(err)
	}

	if user, err := store.Get("maddy"); err == nil {
		log.Println(user)
	}
	if users, err := store.GetAll(); err == nil {
		log.Println(users)
	}
	if err := store.Delete("maddy"); err == nil {
		log.Println("Deleted")
		user, err := store.Get("maddy")
		log.Println(user, err)
	}
}
