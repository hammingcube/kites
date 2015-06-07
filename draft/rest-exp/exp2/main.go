package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"github.com/maddyonline/gits"
	"net/http"
)

type Gist struct {
	Id    string
	Files map[string]File
}

type File struct {
	Content  []byte
	Name     string
	Language string
}

var (
	gist1 = Gist{"abc123", map[string]File{
		"abc.md": File{[]byte("This is amazing"), "abc.md", "Markdown"}, 
		"foo.py": File{[]byte("def func: pass"), "foo.py", "Python"},},
	}
	//gist2 = Gist{"pqr321", []File{{"pqr.md", "Another file", "Markdown"}, {"main.py", "print('hello')", "Python"}}}
	gists = []Gist{gist1}
)

type GistServer struct {
}

func (gs *GistServer) List(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(gists)
}

func DoIt() {
	gs := gits.NewService(&gits.Config{ServerPath: "/Users/mjha/git/data"})
	files := map[string][]byte{
		"abc.txt": []byte("Hello\nHow do you do?\n"),
		"pqr.txt": []byte("New stuff\n"),
		"uvw.txt": []byte("Totally new file\n"),
	}
	gs.AddToRepo("tempting", &gits.User{"123", "maddy"}, files)
}

func (gs *GistServer) Post(w rest.ResponseWriter, r *rest.Request) {
	username := r.PathParam("username")
	log.Println(username)
	DoIt()
	//gist := Gist{}
    gist := map[string]string{}
    err := r.DecodeJsonPayload(&gist)
    if err != nil {
    	log.Println(err)
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    /*log.Printf("%v", gist)
    if gist.Id == "abc123" {
        rest.Error(w, "Gist already exists", 400)
        return
    }*/
    w.WriteJson(gist)
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
	} else {
		w.WriteJson(g)
	}
}

func main() {
	gs := &GistServer{}
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/gists", gs.List),
		rest.Get("/gists/:id", gs.Get),
		rest.Post("/users/:username/gists", gs.Post))
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
