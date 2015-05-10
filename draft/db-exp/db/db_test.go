package db 

import (
	"testing"
	"fmt"
)

func TestMongDB(t *testing.T) {
	db, err := NewConnection("localhost")
	fmt.Printf("%v %v", db, err)
}