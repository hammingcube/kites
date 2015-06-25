package main

import (
"github.com/maddyonline/kites/draft/db-exp/models/store"
"github.com/maddyonline/kites/draft/db-exp/models/gists"
"fmt"
)


func main() {
	DoIt()
}

func DoIt() {
	gstore := gists.NewStore(&store.BoltDBStore{})
	gstore.Open("gists.db", "gists")
	defer gstore.Close()
	g := &gists.Gist{"a", map[string]gists.File{"aa": gists.File{"k", "j"}}}
	gstore.Post(g)
	g2, err := gstore.Get(g.Id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(g)
	fmt.Println(g2)
}
