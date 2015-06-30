package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	_ "github.com/maddyonline/gits"
	"net/http"
	_ "fmt"
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
	cleanup			func()
}

func NewServer() *GistServer {
	DEFAULT_USER := &users.User{
		Id: "default_user", 
		Username: "default_user",
		Email: "default_user@gmail.com",
	}
	
	gstore := gists.NewStore(&store.BoltDBStore{})
	gstore.Open("gists.db")

	ustore := users.NewStore(&store.BoltDBStore{})
	ustore.Open("users.db")

	gs := &GistServer{
		gstore: gstore, 
		ustore: ustore,
		gistsBucket: "gists",
		usersBucket: "users",
		defaultUser: DEFAULT_USER,
		cleanup:	func(){gstore.Close(); ustore.Close();},
	}
	return gs
}


func (gs *GistServer) Get(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	g, err := gs.gstore.Get(gs.gistsBucket, id)
	if u, err := gs.ustore.Get(gs.usersBucket, g.UserId); err == nil {
		bytes, _ := users.ToBytes(u)
		log.Println(string(bytes))
	} else {
		log.Println(err)
	}
	if g == nil || err != nil {
		rest.NotFound(w, r)
		return
	} else {
		w.WriteJson(g)
	}
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
	allGood := func(err error) bool {
		if err != nil {
			log.Println(err)
    		rest.Error(w, err.Error(), http.StatusInternalServerError)
    		return false
    	}
    	return true
	}
	gist := &gists.Gist{}
    err := r.DecodeJsonPayload(gist)
    if !allGood(err) {
    	return
    }
    if gist.UserId == "" {
    	gist.UserId = gs.defaultUser.Id
    }
    owner, err := gs.ustore.Get(gs.usersBucket, gist.UserId)
   
    if !allGood(err) {
    	return
    }

    existing, err := gs.gstore.Get(gs.gistsBucket, gist.Id) 
    if err == store.ErrNotFound {
    	existing = &gists.Gist{}
    } else if !allGood(err) {
    	return
    }

    gists.Update(existing, gist)
    err = gs.gstore.Post(gs.gistsBucket, gist.Id, gist)
    if !allGood(err) {
    	return
    }

    switch owner.Meta {
    case nil:
    	owner.Meta = &users.Data{GistsIds: map[string]string{gist.Id: ""}}
    default:
    	if owner.Meta.GistsIds == nil {
    		owner.Meta.GistsIds = map[string]string{gist.Id: ""}
    	} else {
    		owner.Meta.GistsIds[gist.Id] = ""
    	}
    }
    err = gs.ustore.Post(gs.usersBucket, owner.Id, owner)
    if !allGood(err) {
    	return
    }
    w.WriteJson(existing)
}

func seed(gs *GistServer) {
	g := &gists.Gist{
		Id: "a", 
		Files: map[string]gists.File{"aa": gists.File{
			Name: "k", 
			Language: "j",},
		},
	}
	g.UserId = gs.defaultUser.Id
	err := gs.ustore.Post(gs.usersBucket, gs.defaultUser.Id, gs.defaultUser)
	quit := (err != nil) || gs.gstore.Post(gs.gistsBucket, g.Id, g) != nil
	if quit {
		log.Fatal("Failed to seed")
	}
}

const data = `{"Files": {"aa": {"Content": null, "Name": "k", "Language": "j", "Truncated": false}}, "UserId": "default_user", "Id": "a"}`


func main() {

	g, err := gists.FromBytes([]byte(data))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(g)
	}

	gs := NewServer()
	defer gs.cleanup()
	seed(gs)

	u, err := gs.ustore.Get(gs.usersBucket, gs.defaultUser.Id)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(u)
	log.Println(u.Meta)
	
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
