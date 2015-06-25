package main

import (
	"fmt"
	"github.com/maddyonline/kites/draft/db-exp/models/gists"
	"github.com/maddyonline/kites/draft/db-exp/models/store"
	"github.com/maddyonline/kites/draft/db-exp/models/users"
	"log"
)

func main() {
	DoIt()
}

func UserAction(store users.Store) {
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

func GistAction(gstore gists.Store) {
	g := &gists.Gist{"a", map[string]gists.File{"aa": gists.File{"k", "j"}}}
	gstore.Post(g)
	g2, err := gstore.Get(g.Id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(g)
	fmt.Println(g2)
}

func DoIt() {
	gstore := gists.NewStore(&store.BoltDBStore{})
	gstore.Open("gists.db", "gists")
	defer gstore.Close()

	ustore := users.NewStore(&store.BoltDBStore{})
	ustore.Open("users.db", "users")
	defer ustore.Close()

	GistAction(gstore)
	UserAction(ustore)

}
