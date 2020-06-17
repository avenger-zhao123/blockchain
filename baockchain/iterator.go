package baockchain

import (
	"blockchain/block"
	"errors"
)

//定义一个迭代的结构体类型
type BCIterator struct {
	// 关联需要迭代的对象
	bc *BlockChain
	// 当前迭代的区块hash
	currHash block.Hash
}
//构造迭代
  //确定bc和初始化第一个区块的hash
func NewBCIterator(bc *BlockChain) *BCIterator  {
	a:=&BCIterator{
		bc: bc,
		currHash: bc.lastHash,
	}
	return a
}
//迭代方法，Next();用于迭代和判断迭代是否应该继续
func (bci *BCIterator)Next() (*block.Block,error) {
	// 如果currHash == "" 表示也没有区块了，不用再操作数据库
	if bci.currHash ==""{
		return nil,errors.New("")
	}
	// 当前hash在数据库中，可以获取内容，表示没有错误，可以循环
	block,err :=bci.bc.GetBlock(bci.currHash)
	if err !=nil {
		return nil, err
	}
	// 得到前一个区块的hash
	bci.currHash =block.GetHashPrevBlock()
	// 返回当前区块和没有错误
	return block,nil
}