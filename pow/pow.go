package pow

import (
	"blockchain/block"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)
//建立ProofOfWork 类型结构体
type Proofofwork struct {
	// 需要 pow 工作量区块的区块
	block *block.Block
	// 证明参数目标
	target *big.Int
}
//构造方法
func NewPow(b *block.Block)*Proofofwork {
	p :=&Proofofwork{
		block: b,
		//由于target是指针类型，没有下面的操作就为nil，而nil是不能赋值的
		target:big.NewInt(1),  //new是指声明一个big.Int并且把它传给了target
	}
	//计算target
	p.target.Lsh(p.target,uint(block.HashLen- b.GetBits()+1))
	return p

}
//hash证明(详见test.go)
  // 返回使用的 nonce 和 形成的区块 hash
func (p*Proofofwork) Proof()(int,block.Hash) {
	var hashInt big.Int
	// 基于 block 准备 serviceStr
	serviceStr :=p.block.GenServiceStr()
	// nonce 计数器
	nonce :=0
	fmt.Printf("Target:%d\n",p.target)
	// 迭代计算hash，设置防nonce溢出的条件
	for nonce <= math.MaxInt64 {
		////在哈希之前要将服务器字符串(serviceStr)和随机字符串（nonce计数器）连接
		//q:=serviceStr +strconv.Itoa(nonce)
		// 生成 hash
		hash :=sha256.Sum256([]byte(serviceStr +strconv.Itoa(nonce)))
		// 得到 hash 的 big.Int
		hashInt.SetBytes(hash[:])
		fmt.Printf("Hash  :%s\t%d\n",hashInt.String(), nonce)
		// 判断是否满足难度（数学难题）
		if hashInt.Cmp(p.target) == -1 {
			// 解决问题
			return nonce, block.Hash(fmt.Sprintf("%x",hash))
		}
        nonce++
	}

    return 0,""
}