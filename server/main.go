package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	_ "github.com/maddyonline/gits"
	"net/http"
	"fmt"
	"github.com/maddyonline/kites/models/gists"
	"github.com/maddyonline/kites/models/store"
	"github.com/maddyonline/kites/models/users"
)


type GistServer struct {
	gstore			gists.Store
	ustore			users.Store
	gistsBucket		string
	usersBucket		string
	defaultUser 	*users.User
	userMap			map[string]*users.User
}

func (gs *GistServer) GetAll(w rest.ResponseWriter, r *rest.Request) {
	gists, err := gs.gstore.GetAll(gs.gistsBucket)
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

    if gist.UserId == "" {
    	gist.UserId = gs.defaultUser.Id
    }
    fmt.Printf("HERE!! -->%s<--", gist.UserId)
    fmt.Println(gs.userMap)
    fmt.Printf("Key:%s, Val:%s", gist.UserId, gs.userMap[gist.UserId])

    log.Printf("%v", gist)
    gs.gstore.Post(gs.gistsBucket, gist.Id, gist)
    w.WriteJson(gist)

}

func (gs *GistServer) Get(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	g, err := gs.gstore.Get(gs.gistsBucket, id)
	if g == nil || err != nil {
		rest.NotFound(w, r)
		return
	} else {
		w.WriteJson(g)
	}
}

func UserAction(store users.Store) {
	bucketName := "users"
	u1 := &users.User{
		Id: "maddy", 
		Username: "maddy", 
		Email: "maddy@gmail.com",
	}
	u2 := &users.User{
		Id: "mj", 
		Username: "mj", 
		Email: "mj@gmail.com",
	}
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
	err = store.Post(bucketName, u1.Id, u1)
	if err != nil {
		log.Fatal(err)
	}
}

func Seed(gstore gists.Store, bucketName string) {
	//g := &gists.Gist{"a", map[string]gists.File{"aa": gists.File{"k", "j"}}}
	//gstore.Post(bucketName, g.Id, g)
	//g2, err := gstore.Get(bucketName, g.Id)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(g)
	//fmt.Println(g2)
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

func (gs *GistServer) Initialize() {
	allUsers, _ := gs.ustore.GetAll(gs.usersBucket)
	for i, user := range *allUsers {
		fmt.Printf("Assigning %v to %v\n", user.Id, user)
		gs.userMap[user.Id] = &(*allUsers)[i]
	}

	fmt.Println(gs.userMap)
	fmt.Println("------")

}

func main() {
	DEFAULT_USER := &users.User{
		Id: "default_user", 
		Username: "default_user",
		Email: "default_user@gmail.com",
	}
	
	gstore := gists.NewStore(&store.BoltDBStore{})
	gstore.Open("gists.db")
	defer gstore.Close()

	ustore := users.NewStore(&store.BoltDBStore{})
	ustore.Open("users.db")
	defer ustore.Close()

	ustore.Post("users", DEFAULT_USER.Id, DEFAULT_USER)
	UserAction(ustore)

	gs := &GistServer{
		gstore: gstore, 
		ustore: ustore,
		gistsBucket: "gists",
		usersBucket: "users",
		defaultUser: DEFAULT_USER,
		userMap: map[string]*users.User{},
	}

	gs.Initialize()

	Seed(gs.gstore, gs.gistsBucket)

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
