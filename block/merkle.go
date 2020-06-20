package block

import (
	"crypto/sha256"
	"fmt"
)

//构造默克尔树结构体
type MerkleTree struct {
	//该书的根节点（最上头的hash)
	Root *MerkleNode
}
//构造默克尔数的节点结构体
type MerkleNode struct {
	// 当前节点计算得到的Hash值（可能是交易哈希，也可能是两个子节点的哈希）
	Hash Hash
	// 该节点的子节点。若为叶子节点(交易的hash值），left,right 为 nil
	Left, Right *MerkleNode
}
//构建默克尔树节点（基于是否存在 data，来确定是否为叶子节点）
func NewMerleNode(data []byte,left,right *MerkleNode) *MerkleNode  {   //data 是交易的数据

   hash :=""   //让hash值为空
   //判断是否有data、left、right存在
   if data == nil&&left !=nil&&right !=nil {
   	//不存在data，存在left、right，说明是非叶子节点
   	//hash"左右节点的hash"并转成十六进制输出
   	hash =fmt.Sprintf("%x", sha256.Sum256([]byte(left.Hash +right.Hash)))
   }else {
   	//存在data，说明是叶子节点
   	//hash交易数据（data)并也转成十六进制
	   hash =fmt.Sprintf("%x",sha256.Sum256(data))
	   //为了避免错误，强行让 left和right为空
	   left,right =nil,nil
   }
   //给节点赋值
   s:= &MerkleNode{
   	Hash: hash,
   	Left: left,
   	Right: right,
   }
   return s
}
//构建默克尔树
func NewMerleTree(dataSet [][]byte)*MerkleTree  {   //dataset是交易数据的集合，类型为字节切片的切片
	// 1.构建叶子节点
	nodes :=make([]*MerkleNode,len(dataSet))   //定义叶子节点的集合并确定它的长度
	for i,data :=range dataSet{       //遍历交易数据集合
		 nodes[i] = NewMerleNode(data,nil,nil)  //调用”构建默克尔树节点“函数，做hash运算
	}
	// 2.基于叶子节点，构建上层节点（非叶子节点）
	// 思路从最底层开始，逐层得到节点，外层循环控制的是：层数
	// 循环条件就是，该层中节点的数量大于 1
	levelNodes :=[]*MerkleNode{}   //定义非叶子节点集合
	for len(levelNodes)==0 ||len(levelNodes)>1{  //构建外层循环函数 循环体和选择体都为levelnodes;循环条件是等于0或者大于1
		levelNodes =[]*MerkleNode{}   //清空
		// 基于下一层，构建当前层，下一层的数据，来源于 nodes
		//偶数化
		if len(nodes)%2 ==1{          //奇数个
			nodes =append(nodes,nodes[len(nodes)-1])  // 拷贝最后一个形成偶数个
		}
		//构建内循环函数：每层的个数
		for i,l:=0,len(nodes);i<l;i+=2{
			// 构建当前层
			levelNodes =append(levelNodes,NewMerleNode(nil,nodes[i],
				nodes[i+1]))
		}
		fmt.Println(levelNodes,levelNodes[0].Hash)
		// 将 nodes 更新为当前层的节点
		nodes = levelNodes

	}
	//循环结束，levelNodes 中存在一个节点，就是根节点。
	//构建树即可
	r :=&MerkleTree{
		Root:levelNodes[0],
	}
    return r
}
