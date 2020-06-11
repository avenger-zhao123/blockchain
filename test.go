package main

import (
	"blockchain/baockchain"
)

func main(){
	//区块测试
	//c :=baockchain.NewBlock("","Gensis Block.")
	//fmt.Println(c)
	a :=baockchain.NewBlockchain()
	a.AddGensisBlock()
	a.AddBlock("First Block").AddBlock("Second Block")
	//fmt.Println(a)
	a.Iterate()



}
