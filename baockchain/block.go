package baockchain

import (
	"crypto/sha256"
	"fmt"

	"time"
)

type Hash =string   //设置Hash类型为字符串类型
//创建区块结构体
type Block struct {
	header BlockHeader   //区块头
	txs string           //交易列表
	txCouner int         //交易计数器
	hashCurr Hash        //当前区块的哈希值，算法sha256
}
//创建区块头结构体
type BlockHeader struct {
	version int          //版本信息，节点更新时版本迭代
	hashPrevBlock Hash    //前一个节点的Hash值
	hashMerkRoot Hash     //默克尔树根节点
	time time.Time        //节点生成时间
	bits int              //难度系数
	nonce int             //随机计数器，挖矿相关
}

//设置当前区块hash值
func (b*Block)SetHashCurr() *Block{ //b*Block是方法，SetHashCurr是函数名，*Block是返回值
	//生成头信息的拼接字符串
	headerStr := b.header.Stringify()
	//计算hash值
	b.hashCurr = fmt.Sprintf("%x",sha256.Sum256([]byte(headerStr)))
	return  b
}
//将区块头信息字符串化，利于计算区块头的hash值
func (bh*BlockHeader ) Stringify()string {  //bh*BlockHeader是方法，Stringify()是函数名 string是返回值
	return fmt.Sprintf("%d%s%s%d%d%d",  //返回格式化输出
		bh.version,     //版本信息，节点更新时版本迭代
		bh.hashMerkRoot,   //前一个节点的Hash值
		bh.hashPrevBlock,   //默克尔树根节点
		bh.bits,            //难度系数
		bh.time.UnixNano(), // 得到时间戳，nano 级别
		bh.nonce,         //随机计数器，挖矿相关
		)
}
//构造区块
const nowVersion =0  //设置版本为0
//建立区块
func NewBlock(prevHash Hash,txs string) *Block{
	//实例化区块，指在面向对象的编程中，通常把用类创建对象的过程称为实例化
	b :=&Block{header: BlockHeader{  //区块头
		version: nowVersion,   //版本信息，节点更新时版本迭代
		hashPrevBlock: prevHash,    //默克尔树根节点
		time: time.Now(),      // 得到时间戳，nano 级别
	},
	txs: txs,      //难度系数
	txCouner: 1,}   //交易计数器
	b.SetHashCurr()    //当前区块的哈希值，算法sha256
	return b

}

