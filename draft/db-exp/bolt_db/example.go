package main

import (
    "fmt"
    "log"
    "github.com/boltdb/bolt"
)

var world = []byte("world")

func main() {
    db, err := bolt.Open("bolt.db", 0644, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

   

    for i := 0; i < 100000; i++ {
         key := []byte(fmt.Sprintf("hello%d", i))
         value := []byte("Hello World!")
        // store some data
        err = db.Update(func(tx *bolt.Tx) error {
            bucket, err := tx.CreateBucketIfNotExists(world)
            if err != nil {
                return err
            }

            err = bucket.Put(key, value)
            if err != nil {
                return err
            }
            return nil
        })

        if err != nil {
            log.Fatal(err)
        }
    }
    for i := 0; i < 1000; i++ {
        key := []byte(fmt.Sprintf("hello%d", i))
        // retrieve the data
        err = db.View(func(tx *bolt.Tx) error {
            bucket := tx.Bucket(world)
            if bucket == nil {
                return fmt.Errorf("Bucket %q not found!", world)
            }

            val := bucket.Get([]byte(key))
            fmt.Println(string(val))

            return nil
        })

        if err != nil {
            log.Fatal(err)
        }
    }
}