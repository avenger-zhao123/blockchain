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
		target:big.NewInt(1),
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
		q:=serviceStr +strconv.Itoa(nonce)
		// 生成 hash
		hash :=sha256.Sum256([]byte(q))
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
//验证
  //验证区块已生成的hash值（头信息和随机数）与重新生成的hash值是否相等
  //验证是否满足难题
  //有bool的返回值原因是：该函数的功能就是验证hashcurr是否正确，只有要给主程序返回是否有错！具体错误在哪在主程序体现
func (p *Proofofwork)Validata() bool {
	//重新生成Hash
	serviceStr := p.block.GenServiceStr()
	m := serviceStr + strconv.Itoa(p.block.GetNonce())
	hash := sha256.Sum256([]byte(m))
	//比较
	//由于hash生成是十进制的，而hashcurr是十六进制的，因此要先将hash转化成16进制
	if p.block.GetHashCurr() != fmt.Sprintf("%x", hash) {
		return false
	}
	//比较是否满足难题
	   //还是要有一个目标（要与之前的目标一致）
	   target :=big.NewInt(1)
	   target .Lsh(target,uint(block.HashLen - p.block.GetBits() +1))
	   //设定一个target类型相同的数
	   hashInt :=new(big.Int)
	   //将Hash传到这个数中
	   hashInt.SetBytes(hash[:])
	   //判断大小
	   if hashInt.Cmp(target) !=-1{
	   	return false
	   }
    return true

}
