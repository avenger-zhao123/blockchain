package tx

import (
	"blockchain/wallet"
	"crypto/sha256"
	"fmt"
	"log"
)

// 挖矿奖励金
//单位体系
const satoshi = 1 //1 中本聪
const S = satoshi
const KS = 1000 * S //千
const MS = 1000 * KS // 百万
const GS = 1000 * MS //十亿
const BTC = 100000000 * satoshi
//定义一个挖矿奖励金额
const CoinbaseSubidy =12*BTC

//交易结构体
type TX struct {
	//输入（类型是空结构体）
	Inputs []*Input
	//输出（类型是空结构体）
	Outputs []*Output
	//本次交易hash运算的值
	Hash string

}


//构造交易的方法
func NewTX(ins []*Input,outs []*Output)*TX  {
	//传入输入输出数据
	tx :=&TX{
		Inputs: ins,
		Outputs: outs,
	}
	//调用建立好的交易哈希函数
	tx.SetHash()
	return tx
}

//构造挖矿奖励CoinBase交易
func NewCoinbaseTX(to wallet.Address) *TX  {
	//输入,空的（挖矿奖励只有输出，没有输入）
	ins :=[]*Input{}
	// 输出，仅存在一个输出，给目标为 to 的用户挖矿奖励
	output :=&Output{
		// 常量，存储挖矿奖励金
		Value:CoinbaseSubidy,
	    To:   to,
	}
	//将定义好的输出放到挖矿奖励交易的输出中
	outs :=[]*Output{
		output,
	}
   return NewTX(ins,outs)
	
}
//设置哈希
func (tx *TX)SetHash() *TX {
	//先系列化

	ser ,err := SerializeTX(*tx)
	if err !=nil {
		log.Fatal(err)
	}
	//再对系列化的数据Hash运算
	  //使用hash256计算hash值
	hash :=sha256.Sum256(ser)
	   //将hash值赋值于交易的hash值（由于hash值是十进制的，交易的hash是十六进制因此要将于hash值转化为十六进制）
	tx.Hash =fmt.Sprintf("%x",hash)
	return tx
}

