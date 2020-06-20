package block

import (
	"blockchain/tx"
	"strings"

	//"crypto/sha256"
	"fmt"
	"time"
)

const nowVersion =0  //设置版本为0
type Hash =string   //设置Hash类型为字符串类型
const HashLen = 256  //设置Hash的长度为256
const blockBits =8
//创建区块结构体
type Block struct {
	header   BlockHeader //区块头
	//定义交易结构体后，更新交易列表类型：[]*TX 类型
	txs      []*tx.TX    //交易列表
	txCouner int         //交易计数器
	hashCurr Hash        //当前区块的哈希值，算法sha256
}
//创建区块头结构体
type BlockHeader struct {
	version       int       //版本信息，节点更新时版本迭代
	hashPrevBlock Hash      //前一个节点的Hash值
	hashMerkRoot  Hash      //默克尔树根节点
	time          time.Time //节点生成时间
	bits          int       //难度系数
	nonce         int       //随机计数器，挖矿相关
}
//建立区块
    //在定义交易结构体后，删除交易
func NewBlock(prevHash Hash) *Block {
	//实例化区块，指在面向对象的编程中，通常把用类创建对象的过程称为实例化
	b :=&Block{
		header: BlockHeader{ //区块头
			version:       nowVersion, //版本信息，节点更新时版本迭代
			hashPrevBlock: prevHash,   //默克尔树根节点
			time:          time.Now(), // 得到时间戳，nano 级别
			bits:          blockBits,   //难度系数
		},
		//txs: txs,
		//txCouner: 1,
	}
	//交易计数器
	//b.SetHashCurr()    //当前区块的哈希值，算法sha256
	return b

}
//设置当前区块hash值

//将区块头信息字符串化，利于计算区块头的hash值
//func (bh*BlockHeader ) Stringify()string {  //bh*BlockHeader是方法，Stringify()是函数名 string是返回值
//	return fmt.Sprintf("%d%s%s%d%d%d",  //返回格式化输出
//		bh.version,     //版本信息，节点更新时版本迭代
//		bh.hashMerkRoot,   //默克尔树根节点
//		bh.hashPrevBlock,   //前一个节点的Hash值
//		bh.bits,            //难度系数
//		bh.time.UnixNano(), // 得到时间戳，nano 级别
//		bh.nonce,         //随机计数器，挖矿相关
//		)
//}
////构造区块
//func (b*Block)SetHashCurr() *Block{ //b*Block是方法，SetHashCurr是函数名，*Block是返回值
//	//生成头信息的拼接字符串
//	headerStr := b.header.Stringify()
//	//计算hash值
//	b.hashCurr = fmt.Sprintf("%x",sha256.Sum256([]byte(headerStr)))
//	return  b
//}
//设置当前区块hash值
func ( b*Block)SetHashCurr(hash Hash) *Block {
	//计算hash值
	b.hashCurr =hash
	return b
}
//设置区块的nonce（计数器）
func(b *Block) SetNonce(nonce int)*Block {
	b.header.nonce =nonce
	return b
}
//设置交易
func (b *Block)SetTX(tx *tx.TX) *Block {
	b.txs =append(b.txs,tx)
	b.txCouner ++
	return b
}
//获取当前区块Hash
func ( b*Block)GetHashCurr() Hash {
	return b.hashCurr
}
//获取区块的nonce
func (b *Block) GetNonce() int {
	return b.header.nonce
}
//获取当前区块txs(交易列表）
//func ( b*Block)GetTxs() string {
//	return b.txs
//}
//定义交易结构体后，更新：获取当前区块txs(交易列表)
func ( b*Block)GetTxs() []*tx.TX {
	return b.txs
}
//为了更好的展示区块交易的信息，定义了一个新的获取交易的函数
func (b *Block)GetTxsString() string  {
	//将计数器以16进制的整型打印出来
	show :=fmt.Sprintf("%d tansactions\n",b.txCouner)
	//定义一个空的字符串切片
	txStr :=[]string{}
	//遍历所有的交易信息（如：hash值，输入，输出）
	for  i,t :=range b.txs{
		txStr =append(txStr,fmt.Sprintf("\tindex:%d, Hash: %s, Inputs: %d, Ouputs: %d",
			i, t.Hash,len(t.Inputs),len(t.Outputs)))
		//show +=fmt.Sprintf(" Hash: %s, Inputs: %d, Ouputs: %d",
		//	 t.Hash,len(t.Inputs),len(t.Outputs))
	}
	//返回交易计数器和交易信息
	//return show
	return show + strings.Join(txStr, "\n")

}
//获取当前区块time（生成时间）
func ( b*Block)GetTime() time.Time {
	return b.header.time
}
//获取当前区块HashPrevBlock（前一个节点的Hash值）
func ( b*Block)GetHashPrevBlock() Hash {
	return b.header.hashPrevBlock
}

// 获取bits
func (b *Block) GetBits() int {
	return b.header.bits

}
//获取默克尔树节点
func(b*Block) GetHashMerkRoot(hash Hash)*Block  {
	b.header.hashMerkRoot =hash
	return b
}
// 生成用于 POW（hashCash）的服务字符串
func  (b *Block) GenServiceStr() string {
	return fmt.Sprintf("%d%s%s%d%s", //返回格式化输出
		b.header.version,         //版本信息，节点更新时版本迭代
		b.header.hashMerkRoot,    //前一个节点的Hash值
		b.header.hashPrevBlock,   //默克尔树根节点
		b.header.bits,            //难度系数
		b.header.time.Format("2006-01-02 15:04:05.999999999 -0700 MST"),

	)
}
