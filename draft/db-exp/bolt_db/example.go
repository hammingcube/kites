package main

import (
	"github.com/maddyonline/kites/draft/db-exp/bolt_db/users"
	"github.com/maddyonline/kites/draft/db-exp/bolt_db/gists"
	"log"
)

var (
	gist1 = gists.Gist{"abc123", map[string]gists.File{
		"abc.md": gists.File{[]byte("This is amazing"), "abc.md", "Markdown"}, 
		"foo.py": gists.File{[]byte("def func: pass"), "foo.py", "Python"},
		}, &users.User{"maddy", "maddy", "maddy@gmail.com"},
	}
	all = []gists.Gist{gist1}
)

func main() {
	log.Printf("%s", all)
	log.Printf("%s", all[0].Username)
	log.Printf("%s", all[0].User.Username)

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
