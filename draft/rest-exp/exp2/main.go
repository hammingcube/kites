package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
)

func main() {
	api := rest.NewApi()
	log.Printf("%v", api)
}