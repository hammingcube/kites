package users

import "github.com/boltdb/bolt"
import "encoding/json"
import "fmt"

type BoltDBStore struct {
	dbName     string
	bucketName string
	db         *bolt.DB
}

func (s *BoltDBStore) Open(dbName, bucketName string) error {
	s.dbName = dbName
	s.bucketName = bucketName
	if db, err := bolt.Open(s.dbName, 0644, nil); err == nil {
		s.db = db
		return nil
	} else {
		return err
	}
}

func (s *BoltDBStore) Close() {
	s.db.Close()
}

func (s *BoltDBStore) marshal(u *User) ([]byte, error) {
	return json.Marshal(u)
}

func (s *BoltDBStore) unmarshal(jsonBlob []byte, u *User) error {
	return json.Unmarshal(jsonBlob, u)
}

func (s *BoltDBStore) Post(u *User) error {
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

func (s *BoltDBStore) Delete(key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(s.bucketName)).Delete([]byte(key))
	})
}

func (s *BoltDBStore) Get(key string) (*User, error) {
	user := &User{}
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.bucketName))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", s.bucketName)
		}
		val := bucket.Get([]byte(key))
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

func (s *BoltDBStore) GetAll() ([]User, error) {
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
