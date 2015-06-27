package store

import "github.com/boltdb/bolt"
import "fmt"

type BoltDBStore struct {
	dbName     string
	db         *bolt.DB
}

func (s *BoltDBStore) Open(dbName string) error {
	s.dbName = dbName
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

func (s *BoltDBStore) Post(bucketName, key, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
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

func (s *BoltDBStore) Delete(bucketName, key []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Delete(key)
	})
}

func (s *BoltDBStore) Get(bucketName, key []byte) ([]byte, error) {
	var value []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucketName)
		}
		value = bucket.Get(key)
		if value == nil {
			return fmt.Errorf("User %q not found!", key)
		}
		return nil
	})
	return value, err
}

func (s *BoltDBStore) GetAll(bucketName []byte) ([][]byte, error) {
	values := [][]byte{}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			values = append(values, v)
		}
		return nil
	})
	return values, err
}
