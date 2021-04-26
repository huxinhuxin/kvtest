package main

import (
	"errors"

	"github.com/cockroachdb/pebble"
)

type Pebble struct {
	DB     *pebble.DB
	Getnum uint64
	Setnum uint64
	Delnum uint64
	Bsync  bool
}

func NewPebble() DBInterface {
	return &Pebble{}
}

func (db *Pebble) Open(path string, sync bool) error {

	database, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return err
	}
	db.DB = database
	db.Bsync = sync
	return nil
}

func (db *Pebble) Close() error {
	return db.DB.Close()
}

func (db *Pebble) Get(key []byte) ([]byte, error) {
	db.Getnum++
	v, closer, err := db.DB.Get(key)

	if v == nil {
		err = errors.New("keyNotFound")
	}
	if closer != nil {
		err = closer.Close()
	}
	return v, err
}

func (db *Pebble) Set(key, val []byte) error {
	opts := &pebble.WriteOptions{Sync: db.Bsync}
	db.Setnum++
	return db.DB.Set(key, val, opts)
}

func (db *Pebble) Del(key []byte) error {
	db.Delnum++
	return db.DB.Delete(key, &pebble.WriteOptions{Sync: db.Bsync})
}

func (db *Pebble) GetAll() (int, error) {
	var cout int
	iter := db.DB.NewIter(&pebble.IterOptions{})
	defer iter.Close()
	iter.First()

	for iter.First(); iter.Valid(); iter.Next() {
		cout++
	}
	return cout, nil
}
