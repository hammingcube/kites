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
	bucketName := "users"
	u1 := &users.User{"maddy", "maddy", "maddy@gmail.com"}
	u2 := &users.User{"mj", "mj", "mj@gmail.com"}
	err := store.Post(bucketName, u1.Id, u1)
	err = store.Post(bucketName, u2.Id, u2)
	if err != nil {
		log.Fatal(err)
	}

	if user, err := store.Get(bucketName, "maddy"); err == nil {
		log.Println(user)
	}
	if users, err := store.GetAll(bucketName); err == nil {
		log.Println(users)
	}
	if err := store.Delete(bucketName, "maddy"); err == nil {
		log.Println("Deleted")
		user, err := store.Get(bucketName, "maddy")
		log.Println(user, err)
	}
}

func GistAction(gstore gists.Store) {
	bucketName := "gists"
	g := &gists.Gist{"a", map[string]gists.File{"aa": gists.File{"k", "j"}}}
	gstore.Post(bucketName, g.Id, g)
	g2, err := gstore.Get(bucketName, g.Id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(g)
	fmt.Println(g2)
}

func DoIt() {
	gstore := gists.NewStore(&store.BoltDBStore{})
	gstore.Open("gists.db")
	defer gstore.Close()

	ustore := users.NewStore(&store.BoltDBStore{})
	ustore.Open("users.db")
	defer ustore.Close()

	GistAction(gstore)
	UserAction(ustore)

}
