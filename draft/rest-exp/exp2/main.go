package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"log"
)

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	app := rest.AppSimple(func(w rest.ResponseWriter, r *rest.Request){
			w.WriteJson(map[string]string{"Body": "Hello World!"})
		})
	api.SetApp(app)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}