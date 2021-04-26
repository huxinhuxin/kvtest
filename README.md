# kvtest
这是一个简单的测试框架，测试几个kv数据库的性能，现在包括boltdb,rocksdb,pebbledb,badgerdb,rocksdb只支持linux，并且需要先在机器上安装rocksdb。

启动参数：
-f 指定存放目录，默认为当前目录的mytestdata
-d true或false,指定是否删除原来的目录，默认为true
-a 指定运行次数，默认为1000

设定假如设定a次，将会依次执行后续操作：
1.	顺序a次写
2.	随机读a次
3.	随机删a次
4.	随机写a次
5.	遍历所有Key,

-sync  true或false ,表示是否同步写盘，默认为true
-t  指定测试的kv类型，现在只有bolt,rocks,pebble,badger四种，默认为全部。
