package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"log"
)

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get("/gists", func(w rest.ResponseWriter, r *rest.Request) {
				a := []map[string][]string{
					{"files": {"file1", "file2"}},
					{"files": {"file3", "file4"}}}
				w.WriteJson(a)
				}))
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}