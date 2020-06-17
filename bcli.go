package main

import (
	"blockchain/baockchain"
	"blockchain/wallet"
	"flag"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"strings"
)

func main()  {
	// 初始化数据库
	  //数据库连接
	d :="data"
	db,err :=leveldb.OpenFile(d,nil)
	if err !=nil{
		log.Fatal(err)
	}
	  //释放数据库连接
	defer db.Close()
	// 初始化区块链
	a :=baockchain.NewBlockchain(db)
	////添加创世区块
	//a.AddGensisBlock()
	//声明一个变量
    var arg1 string
	// 若用户指定了参数，则第一个用户参数为命令参数,命令至少有两个
	//命令至少有两个原因：os.Args 用于获取命令行上的全部参数，为 []string 结构，其中 0 元素固定，为当前执行的脚本名。
      if len(os.Args) >= 2 {
      	//当命令至少为两个时，将os.Args 下标为1赋予给arg1
      	arg1 =os.Args[1]

	  }else {
	  	//当命令参数少于2个时，arg1变量为空
		  arg1 =""
	  }
	// 基于命令参数，执行对应的功能
	  switch strings.ToLower(arg1) {
	  //添加一个选择：增加区块命令
		 // 实现的效果，执行如下两个命令：
	  	     //>go run bcli.go createblock "A send 1 to B"
	  	     //>go run bcli.go createblock -txs="A send 1 to B"
	  	//我们选择性实现，增加通用性，使用 -txs 的模式。（有-txs后面可以有好几个参数，没有就只能跟一个参数
	  	//完成对 -txs 命令行标志 flag 的解析，解析 flag go 提供了flag包完成
	  	case "create:block":
		  // 为 createblock 命令增加一个 flag 集合。标志集合
		  //flag为“-txs"的参数，flag 集合就是多个“-txs"的参数
		  //下面参数1为flag 集合的名字，与选择名字一样，易于识别
		  f :=flag.NewFlagSet("create:block",flag.ExitOnError)
			// 在集合中，添加需要解析的 flag 标志
			//参数的名字为"txs",参数的默认值为空，参数的帮助信息为空
		  address :=f.String("address","","")
		  // 解析命令行参数,命令参数[0]是当前执行的脚本名，参数[1]是-txs
		  f.Parse(os.Args[2:])
		  //判断是否解析成功
		  //if !f.Parsed() {
		  //	log.Fatal("createblock args parsed error")
		  //}
		  //fmt.Println(txs,*txs)
		  //完成区块的创建
		  a.AddBlock(*address)
		  //展示全部区块
	  case "show":
	  	a.Iterate()
	  case "init":
	  	fs :=flag.NewFlagSet("init",flag.ExitOnError)
	  	address :=fs.String("address","","")
	  	fs.Parse(os.Args[2:])
		  // 清空
	  	a.Clear()
		  // 添加创世区块
	  	a.AddGensisBlock(*address)
	  case "create:wallet":
		  // 命令行标志集（参数集 -flag）
	  	fs :=flag.NewFlagSet("create:wallet",flag.ExitOnError)
		  // pass 标志, *string
	  	pass :=fs.String("pass","","")
	  	w :=wallet.NewWallet(*pass)
	  	fmt.Printf("you mnemonic: %s\n",w.GetMnemonic())
	  	fmt.Printf("you address: %s \n",w.Address)
	  case "balance":
	  	fs :=flag.NewFlagSet("balance",flag.ExitOnError)
	  	address :=fs.String("address","","")
	  	fs.Parse(os.Args[2:])

	  	fmt.Printf("Address:%s\nBalance:%d\n",
	  		*address,a.GetBalance(*address),
	  		)



	  case "help":
		  fallthrough  //贯穿
	  default:   //预设
		  Usage()
	}
}
func Usage()  {
	fmt.Println("bcli is a tool for Blockchain.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("\t%s\t%s\n", "bcli create:block -txs=<txs>", "create block on blockchain.")
	fmt.Printf("\t%s\t%s\n", "bcli create:wallet -pass=<pass>", "create wallet base on pass.")
	fmt.Printf("\t%s\t%s\n", "bcli init -address=<address>", "initial blockchain")
	fmt.Printf("\t%s\t%s\n", "bcli balance -address=<address>", "get address 's balance")
	fmt.Printf("\t%s\t\t\t%s\n", "bcli help", "help info for bcli")
	fmt.Printf("\t%s\t\t\t%s\n", "bcli show", "show blocks in chain.")

}