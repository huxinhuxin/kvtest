package main

type Gorocksdb struct {
	//DB     *gorocksdb.DB
	Getnum uint64
	Setnum uint64
	Delnum uint64
	Bsync  bool
}

func NewRocksDB() DBInterface {
	return &Gorocksdb{}
}

func (db *Gorocksdb) Open(path string, sync bool) error {
	// opt := gorocksdb.NewDefaultOptions()
	// opt.SetCreateIfMissing(true)

	// database, err := gorocksdb.OpenDb(opt, path)
	// if err != nil {
	// 	return err
	// }
	// db.DB = database
	// db.Bsync = sync
	return nil
}

func (db *Gorocksdb) Close() error {
	//db.DB.Close()
	return nil
}

func (db *Gorocksdb) Get(key []byte) ([]byte, error) {
	// db.Getnum++
	// opt := gorocksdb.NewDefaultReadOptions()

	// v, err := db.DB.Get(opt, key)

	// if v == nil {
	// 	err = errors.New("keyNotFound")
	// }
	// return v.Data(), err
	return nil, nil
}

func (db *Gorocksdb) Set(key, val []byte) error {
	// opts := gorocksdb.NewDefaultWriteOptions()
	// db.Setnum++
	// opts.SetSync(db.Bsync)
	// return db.DB.Put(opts, key, val)
	return nil
}

func (db *Gorocksdb) Del(key []byte) error {
	// db.Delnum++
	// opts := gorocksdb.NewDefaultWriteOptions()
	// opts.SetSync(db.Bsync)
	// return db.DB.Delete(opts, key)
	return nil
}

func (db *Gorocksdb) GetAll() (int, error) {
	var cout int
	// opt := gorocksdb.NewDefaultReadOptions()
	// iter := db.DB.NewIterator(opt)
	// defer iter.Close()
	// for iter.SeekToFirst(); iter.Valid(); iter.Next() {
	// 	cout++
	// }
	return cout, nil
}
