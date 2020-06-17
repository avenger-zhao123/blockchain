package baockchain

import (
	"blockchain/wallet"
	"fmt"
)

//统计余额
func (bc *BlockChain)GetBalance(address wallet.Address) int   {
	//将余额初始化
	 balance :=0
	 //遍历属于这个address的未花费的输出
	 for _,utxo:=range bc.FindUTXO(address){
	 	//将这些输出的值加入余额（由于余额是int型，而append中的参数是字节型切片，不能用append）
	 	balance +=utxo.Value
	 }
	 fmt.Println(balance)
	 return balance

}
