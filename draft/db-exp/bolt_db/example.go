package main

import (
    "fmt"
    "log"
    "github.com/boltdb/bolt"
    "encoding/json"
)

type User struct {
    Id          string
    Username    string
    Email       string
}

type UsersStore struct {
    dbName      string
    bucketName  string
    db          *bolt.DB
}

type Store interface {
    Open(dbName, bucketName string) error
    Close()
    Get(key string) (*User, error)
    GetAll() ([]User, error)
    Post(u *User) error
    Delete(key string) error
}

func (s *UsersStore) Open(dbName, bucketName string) error {
    s.dbName = dbName
    s.bucketName = bucketName
    if db, err := bolt.Open(s.dbName, 0644, nil); err == nil {
        s.db = db
        return nil
    } else {
        return err
    }
}

func (s *UsersStore) Close() {
    s.db.Close()
}

func (s *UsersStore) marshal(u *User) ([]byte, error) {
    return json.Marshal(u)
}

func (s *UsersStore) unmarshal(jsonBlob []byte, u *User) error {
    return json.Unmarshal(jsonBlob, u)
}

func (s *UsersStore) Post(u *User) error {
    return s.db.Update(func(tx *bolt.Tx) error {
        key := []byte(u.Id)
        value, err := s.marshal(u)
        if err != nil {
            return err
        }
        bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucketName))
        if err != nil {
            return err
        }
        err = bucket.Put(key, value)
        if err != nil {
                return err
        }
        return nil
    })
}

func (s *UsersStore) Delete(key string) error {
    return s.db.Update(func(tx *bolt.Tx) error {
        return tx.Bucket([]byte(s.bucketName)).Delete([]byte(key))
    })
}


func (s *UsersStore) Get(key string) (*User, error) {
    user := &User{}
    err := s.db.View(func(tx *bolt.Tx) error {
            bucket := tx.Bucket([]byte(s.bucketName))
            if bucket == nil {
                return fmt.Errorf("Bucket %q not found!", s.bucketName)
            }
            val := bucket.Get([]byte(key)); 
            if val == nil {
                return fmt.Errorf("User %q not found!", key)
            }
            if err := s.unmarshal(val, user); err != nil {
                return err
            } 
            return nil
        })
    return user, err
}

func (s *UsersStore) GetAll() ([]User, error) {
    users := []User{}
    err := s.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(s.bucketName))
        c := b.Cursor()
        for k, v := c.First(); k != nil; k, v = c.Next() {
            u := &User{}
            if err := s.unmarshal(v, u); err != nil {
                return err
            }
            users = append(users, *u)
        }
        return nil
    })
    return users, err
}


func main() {
    store := &UsersStore{}
    DoIt(store)
}

func DoIt(store Store) {
    if err := store.Open("usersdb", "users"); err != nil {
        log.Fatal(err)
    }
    defer store.Close()
    u1 := &User{"maddy", "maddy", "maddy@gmail.com"}
    u2 := &User{"mj", "mj", "mj@gmail.com"}
    err := store.Post(u1)
    err = store.Post(u2)
    if err != nil {
        log.Fatal(err)
    }
    
    if user, err := store.Get("maddy"); err == nil {
        log.Println(user)
    }
    if users, err := store.GetAll(); err == nil {
        log.Println(users)
    }
    if err := store.Delete("maddy"); err == nil {
        log.Println("Deleted")
        user, err := store.Get("maddy")
        log.Println(user, err)
    }
}