package store

import "github.com/boltdb/bolt"
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

func (s *BoltDBStore) Post(key, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
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

func (s *BoltDBStore) Delete(key []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(s.bucketName)).Delete(key)
	})
}

func (s *BoltDBStore) Get(key []byte) ([]byte, error) {
	var value []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.bucketName))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", s.bucketName)
		}
		value = bucket.Get(key)
		if value == nil {
			return fmt.Errorf("User %q not found!", key)
		}
		return nil
	})
	return value, err
}

func (s *BoltDBStore) GetAll() ([][]byte, error) {
	values := [][]byte{}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			values = append(values, v)
		}
		return nil
	})
	return values, err
}
