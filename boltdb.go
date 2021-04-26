package main

import (
	"errors"

	"github.com/boltdb/bolt"
)

type Boltdb struct {
	DB     *bolt.DB
	Getnum uint64
	Setnum uint64
	Delnum uint64
	Bsync  bool
}

var bucketName []byte

func NewBoltDB() DBInterface {
	return &Boltdb{}
}

func (db *Boltdb) Open(path string, sync bool) error {

	database, err := bolt.Open(path, 0600, &bolt.Options{})
	if err != nil {
		return err
	}
	db.DB = database
	bucketName = S2b("mybucket")
	db.DB.NoSync = !sync
	return nil
}

func (db *Boltdb) Close() error {
	return db.DB.Close()
}

func (db *Boltdb) Get(key []byte) ([]byte, error) {
	var ret []byte
	err := db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return errors.New("bukey not fond")
		}
		db.Getnum++
		v := b.Get(key)
		if v == nil {
			return errors.New("keyNotFound")
		}
		ret = make([]byte, len(v))
		copy(ret, v)
		return nil
	})
	return ret, err
}

func (db *Boltdb) Set(key, val []byte) error {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			var err error
			b, err = tx.CreateBucket(bucketName)
			if err != nil {
				return err
			}
		}
		db.Setnum++
		return b.Put(key, val)
	})
	return err
}

func (db *Boltdb) Del(key []byte) error {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b != nil {
			db.Delnum++
			return b.Delete(key)
		}
		return nil
	})
	return err
}

func (db *Boltdb) GetAll() (int, error) {
	var cout int
	err := db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b != nil {
			c := b.Cursor()

			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				//fmt.Printf("key=%s, value=%s\n", k, v)
				cout++
				//v[10] = 15
			}
		}

		return nil
	})
	return cout, err
}
