package main

import "fmt"

type Gist struct {
	Id    string
	Files []File
}

type File struct {
	Content  string
	Language string
}

var (
	gist1 = Gist{"abc123", []File{{"This is amazing", "md"}, {"def func: pass", "py"}}}
	gist2 = Gist{"pqr321", []File{{"Another file", "md"}, {"print('hello')", "py"}}}
	fixture = []Gist{gist1, gist2}
)

func main() {
	fmt.Println(fixture)
}