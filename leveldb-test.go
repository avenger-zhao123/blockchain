package main

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)
func main() {
	//打开数据库连接：（打开文件或目录句柄）
	dbpath :="testdb"  //数据库目录
	db,err :=leveldb.OpenFile(dbpath,nil)//若testdb不存在，系统会自动维护
	if err != nil {
		log.Fatal(err)  ///打不开的原因是go语言在文件系统中没有权限
	}
	//设置（创建，更新）一个key（leveldb是key-value型数据库）
	  key :="avenger zhao"  //初始化key
	  if err :=db.Put([]byte(key),[]byte("premier league"),nil);err!=nil{
	  	log.Fatal(err)
	  }
	  log.Println("put success")
	//读取key
	data,err :=db.Get([]byte(key),nil)
	if err !=nil{
		log.Fatal(err)
	}
	fmt.Println(data ,string(data))
}
