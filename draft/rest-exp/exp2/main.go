package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
)

type Gist struct {
	Id    string
	Files []File
}

type File struct {
	Name     string
	Content  string
	Language string
}

var (
	gist1 = Gist{"abc123", []File{{"abc.md", "This is amazing", "Markdown"}, {"foo.py", "def func: pass", "Python"}}}
	gist2 = Gist{"pqr321", []File{{"pqr.md", "Another file", "Markdown"}, {"main.py", "print('hello')", "Python"}}}
	gists = []Gist{gist1, gist2}
)

type GistServer struct {
}

func (gs *GistServer) List(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(gists)
}

func (gs *GistServer) Get(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	var g *Gist
	for _, gist := range gists {
		if gist.Id == id {
			g = &gist
			break
		}
	}
	if g == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(g)
}

func main() {
	gs := &GistServer{}
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/gists", gs.List),
		rest.Get("/gists/:id", gs.Get))
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
