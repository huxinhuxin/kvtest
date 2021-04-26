package main

import (
	"github.com/dgraph-io/badger/v3"
)

type Badgerdb struct {
	DB     *badger.DB
	Getnum uint64
	Setnum uint64
	Delnum uint64
	Bsync  bool
}

func NewBadgerDB() DBInterface {
	return &Badgerdb{}
}

func (db *Badgerdb) Open(path string, sync bool) error {

	opt := badger.DefaultOptions(path).WithSyncWrites(sync)
	//opt.SyncWrites = sync
	database, err := badger.Open(opt)
	if err != nil {
		return err
	}
	db.DB = database
	return nil
}

func (db *Badgerdb) Close() error {
	return db.DB.Close()
}

func (db *Badgerdb) Get(key []byte) ([]byte, error) {
	var ret []byte
	err := db.DB.View(func(tx *badger.Txn) error {

		db.Getnum++
		v, err := tx.Get(key)
		if err != nil {
			return err
		}
		ret, err = v.ValueCopy(nil)
		return nil
	})
	return ret, err
}

func (db *Badgerdb) Set(key, val []byte) error {
	err := db.DB.Update(func(tx *badger.Txn) error {
		db.Setnum++
		return tx.Set(key, val)
	})
	return err
}

func (db *Badgerdb) Del(key []byte) error {
	err := db.DB.Update(func(tx *badger.Txn) error {

		return tx.Delete(key)
	})
	return err
}

func (db *Badgerdb) GetAll() (int, error) {
	var cout int
	err := db.DB.View(func(tx *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		it := tx.NewIterator(opt)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			// item := it.Item()
			// k := item.Key()
			// err := item.Value(func(v []byte) error {
			//   fmt.Printf("key=%s, value=%s\n", k, v)
			//   return nil
			// })
			// if err != nil {
			//   return err
			// }
			cout++
		}

		return nil
	})
	return cout, err
}
