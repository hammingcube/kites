package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	_ "github.com/maddyonline/gits"
	"net/http"
	"fmt"
	"github.com/maddyonline/kites/draft/db-exp/models/gists"
	"github.com/maddyonline/kites/draft/db-exp/models/store"
	"github.com/maddyonline/kites/draft/db-exp/models/users"
)


type GistServer struct {
	gstore	gists.Store
	bucketName	string
}

func (gs *GistServer) GetAll(w rest.ResponseWriter, r *rest.Request) {
	gists, err := gs.gstore.GetAll(gs.bucketName)
	if gists == nil || err != nil {
		rest.NotFound(w, r)
		return
	} else {
		w.WriteJson(gists)
	}
}

func (gs *GistServer) Post(w rest.ResponseWriter, r *rest.Request) {
	gist := &gists.Gist{}
    err := r.DecodeJsonPayload(gist)
    if err != nil {
    	log.Println(err)
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    log.Printf("%v", gist)
    gs.gstore.Post(gs.bucketName, gist.Id, gist)
    w.WriteJson(gist)

}

func (gs *GistServer) Get(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	g, err := gs.gstore.Get(gs.bucketName, id)
	if g == nil || err != nil {
		rest.NotFound(w, r)
		return
	} else {
		w.WriteJson(g)
	}
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

func Seed(gstore gists.Store, bucketName string) {
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

	Seed(gstore, "gists")
	UserAction(ustore)

}

func main() {
	gstore := gists.NewStore(&store.BoltDBStore{})
	gstore.Open("gists.db")
	defer gstore.Close()
	gs := &GistServer{gstore: gstore, bucketName: "gists"}

	Seed(gs.gstore, gs.bucketName)

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/gists", gs.GetAll),
		rest.Get("/gists/:id", gs.Get),
		rest.Post("/gists", gs.Post))
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
