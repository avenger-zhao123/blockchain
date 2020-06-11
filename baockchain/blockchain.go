package baockchain

import (
	"fmt"
	"time"
)

//建立区块链类型结构体
type BlockChain struct {
	lastHash Hash           //最后一个区块的哈希
	blocks map[Hash]*Block  //全部区块信息，由区块哈希作为key来检索
}
//建立区块链
func NewBlockchain()*BlockChain {
	//实例化区块，指在面向对象的编程中，通常把用类创建对象的过程称为实例化
	bc := &BlockChain{
		blocks : map[Hash]*Block{},  //全部区块信息，由区块哈希作为key来检索
	}
	return bc
}
//添加区块到区块链上
//bc *BlockChain是方法，AddBlock是函数名， 提供区块的数据，目前是字符串，*BlockChain是返回值
func (bc *BlockChain) AddBlock(txs string) *BlockChain {
	//构建区块
	b :=NewBlock(bc.lastHash,txs)
	//将区块加入到链的存储结构中
	bc.blocks[b.hashCurr] =b
	//将最后的区块哈希设置为当前区块
	bc.lastHash =b.hashCurr
	return bc

}
//添加创世区块（第一个区块）
//bc *BlockChain是方法，AddGensisBlock是函数名，*BlockChain是返回值
func (bc *BlockChain)AddGensisBlock() *BlockChain  {
	//校验是否可以添加创世区块
	if bc.lastHash !="" {
		//已经存在区块，不需要再添加创世区块
		return bc
	}
	return bc.AddBlock("Founding block")
}
//迭代展示区块的方法 （方便之后的测试）
func (bc *BlockChain) Iterate() {
	//通过for循环遍历出当前区块
	for hash :=bc.lastHash;hash !=""; {
		//b作为blocks的下标
		b :=bc.blocks[hash]
		//打印区块的Hash值
		fmt.Println("HashCurr:", b.hashCurr)
		//打印区块的交易列表
		fmt.Println("Txs", b.txs)
		//打印节点生成时间
		fmt.Println("Time", b.header.time.Format(time.UnixDate))
		//打印前一个节点的哈希值
		fmt.Println("HashPrev",b.header.hashPrevBlock)
		hash =b.header.hashPrevBlock

	}
}