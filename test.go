package main

import (
	"blockchain/baockchain"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

func main(){
	//区块测试
	//c :=baockchain.NewBlock("","Gensis Block.")
	//fmt.Println(c)
	//数据库连接
	d :="data"
	db,err :=leveldb.OpenFile(d,nil)
	if err !=nil{
		log.Fatal(err)
	}
	//释放数据库连接
	defer db.Close()

    //区块链测试
	a :=baockchain.NewBlockchain(db)
	a.AddGensisBlock() // ->
	a.AddBlock("First Block").AddBlock("Second Block")
	//fmt.Println(a)
	a.Iterate()



}
