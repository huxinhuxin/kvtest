package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type DBInterface interface {
	Set(key, val []byte) error
	Get(key []byte) ([]byte, error)
	Del(key []byte) error
	Open(path string, sync bool) error
	GetAll() (int, error)
	Close() error
}

type timeStatistics struct {
	sequewrite time.Duration
	randRead   time.Duration
	randDel    time.Duration
	randWrite  time.Duration
	getAll     time.Duration
	total      time.Duration
}

func (t timeStatistics) String() string {
	return fmt.Sprintf("\n%-30s%v\n%-30s%v\n%-30s%v\n%-30s%v\n%-30s%v\n%-30s%v\n",
		"seque write cost time:", t.sequewrite,
		"rand read cost time:", t.randRead,
		"rand del cost time:", t.randDel,
		"rand write cost:", t.randWrite,
		"getall cost:", t.getAll,
		"total cost:", t.total)
}

var g_count int

const prikey = "user/"

var g_cache []byte

var g_opt int

func sequencewrite(db DBInterface, data []byte) time.Duration {
	log.Println("sequence write begin")
	t1 := time.Now()
	for i := 0; i < g_count; i++ {
		key := append([]byte(prikey), I2b(uint64(i))...)
		db.Set(key, data)
	}
	t2 := time.Since(t1)
	log.Println("sequence write end, cost time ", t2)

	return t2
}

func randwrite(db DBInterface, data []byte) time.Duration {
	log.Println("randwrite  begin")
	t1 := time.Now()
	for i := 0; i < g_count; i++ {
		tmp := rand.Int31n(int32(g_count))
		key := append([]byte(prikey), I2b(uint64(tmp))...)
		db.Set(key, data)
	}
	t2 := time.Since(t1)
	log.Println("randwrite  end, cost time ", t2)
	return t2
}

func randRead(db DBInterface) time.Duration {
	log.Println("randRead  begin")
	t1 := time.Now()
	count := 0
	for i := 0; i < g_count; i++ {
		tmp := rand.Int31n(int32(g_count))
		key := append([]byte(prikey), I2b(uint64(tmp))...)
		v, err := db.Get(key)
		if err != nil {
			count++
		}
		copy(g_cache, v)

	}
	log.Println("randRead err count = ", count)
	t2 := time.Since(t1)
	//log.Println("randRead  end, cost time ", t2)
	return t2
}

func randdel(db DBInterface) time.Duration {
	log.Println("randdel  begin")
	t1 := time.Now()
	count := 0
	for i := 0; i < g_count; i++ {
		tmp := rand.Int31n(int32(g_count))
		key := append([]byte(prikey), I2b(uint64(tmp))...)
		err := db.Del(key)
		if err != nil {
			count++
		}

	}
	t2 := time.Since(t1)
	log.Println("randdel err count = ", count)
	log.Println("randdel  end, cost time ", t2)
	return t2
}

func readall(db DBInterface) time.Duration {
	log.Println("readall  begin")
	t1 := time.Now()
	cout, _ := db.GetAll()
	t2 := time.Since(t1)
	log.Println("readall get count = ", cout)
	log.Println("readall  end, cost time ", t2)
	return t2
}

func createtestdb(dbtype, path string, bsync bool) *timeStatistics {
	var db DBInterface
	newpath := path
	switch dbtype {
	case "bolt":
		db = NewBoltDB()
		newpath += "/bolt.db"
	case "pebble":
		db = NewPebble()
		newpath += "/pebble"
	case "badger":
		db = NewBadgerDB()
		newpath += "/badger"
	case "rocks":
		db = NewRocksDB()
		newpath += "/rocksdb"
	default:
		log.Printf("%s dbtype can't test\n", dbtype)
		return nil
	}

	log.Printf("begin test %s !!!!!!!!!!!!!\n", dbtype)
	ti := comtestdb(db, newpath, bsync)
	log.Printf("end test %s !!!!!!!!!!!!!\n", dbtype)
	return ti
}

func comtestdb(db DBInterface, path string, bsync bool) *timeStatistics {
	t1 := time.Now()
	if err := db.Open(path, bsync); err != nil {
		fmt.Println("open err")
		return nil
	}
	writedata := randomString(4096)

	timeStat := &timeStatistics{}

	switch g_opt {
	case 1:
		timeStat.sequewrite = sequencewrite(db, writedata)
		timeStat.randWrite = randwrite(db, writedata)
	case 2:
		timeStat.randDel = randdel(db)
	case 3:
		timeStat.randRead = randRead(db)
		timeStat.getAll = readall(db)
	case 4:
		timeStat.sequewrite = sequencewrite(db, writedata)
		timeStat.randDel = randdel(db)
		timeStat.randWrite = randwrite(db, writedata)

	default:
		timeStat.sequewrite = sequencewrite(db, writedata)
		timeStat.randRead = randRead(db)
		timeStat.randDel = randdel(db)
		timeStat.randWrite = randwrite(db, writedata)
		timeStat.getAll = readall(db)
	}

	db.Close()
	timeStat.total = time.Since(t1)
	log.Printf("%d counts sequence write, rand read,rand del,rand write,get all\n", g_count)
	log.Println("db close ,total cost time ", timeStat.total)
	return timeStat
}
func main() {
	dFile := flag.String("f", "mytestdata", "db file path")
	aNum := flag.Int("a", 1000, "want to add article num")
	bsync := flag.Bool("sync", true, "b sync")
	delolddb := flag.Bool("d", true, "delete old")
	dbtype := flag.String("t", "all", "bolt,pebble,badger,all")
	parse := flag.Int("p", 0, "use parse")
	opt := flag.Int("o", 0, " 0:all, 1:only write  2:only delete 3:only read")
	flag.Parse()
	g_opt = *opt
	if *parse == 1 {
		fmt.Println("parse file")
		disk := parsedisk()
		if disk != nil {
			printdisk(disk)
		}

		cpu := parsecpu()
		if cpu != nil {
			printcpu(cpu)
		}

		diskx := parsedisk1()
		if diskx != nil {
			printdisk1(diskx)
		}
		return
	}
	g_cache = make([]byte, 4096)
	g_count = *aNum
	fmt.Println(os.Args)
	mapstatis := make(map[string]*timeStatistics)

	if *delolddb {
		os.RemoveAll(*dFile)
	}
	os.Mkdir(*dFile, 0666)
	os.RemoveAll("data")
	os.Mkdir("data", 0666)

	//启动iostat
	list1 = make([]context.CancelFunc, 0)
	wg := &sync.WaitGroup{}
	go startproc([]string{"-d", "-k", "1"}, "data/disk", wg)
	go startproc([]string{"-c", "1"}, "data/cpu", wg)
	go startproc([]string{"-d", "-x", "1"}, "data/diskx", wg)
	switch *dbtype {
	case "all":
		t1 := createtestdb("bolt", *dFile, *bsync)
		worktomap(mapstatis, "bolt", t1)
		t2 := createtestdb("pebble", *dFile, *bsync)
		worktomap(mapstatis, "pebble", t2)
		t3 := createtestdb("badger", *dFile, *bsync)
		worktomap(mapstatis, "badger", t3)
		t4 := createtestdb("rocks", *dFile, *bsync)
		worktomap(mapstatis, "rocks", t4)
	default:
		t := createtestdb(*dbtype, *dFile, *bsync)
		worktomap(mapstatis, *dbtype, t)
	}

	for key, val := range mapstatis {
		fmt.Println("---------------", key, "--------------")
		fmt.Println(val)
		fmt.Println("-----------------------------------")
	}

	for _, cacel := range list1 {
		cacel()
	}

	wg.Wait()

}

func worktomap(m map[string]*timeStatistics, key string, t *timeStatistics) {
	if t == nil {
		return
	}
	if val, ok := m[key]; ok {

		val.sequewrite = (val.sequewrite + t.sequewrite) / 2
		val.randRead = (val.randRead + t.randRead) / 2
		val.randDel = (val.randDel + t.randDel) / 2
		val.randWrite = (val.randWrite + t.randWrite) / 2
		val.getAll = (val.getAll + t.getAll) / 2
		val.total = (val.total + t.total) / 2

	} else {
		m[key] = t
	}
}
