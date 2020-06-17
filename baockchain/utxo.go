package baockchain

import (
	"blockchain/tx"
	"blockchain/wallet"
)
// 找到哪些属于address的未被花费的输出Unspend
   //UTXO：没有作为其他交易输入的输出，称为为花费的交易输出。
func (bc *BlockChain)FindUTXO(address wallet.Address) []*tx.Output  {
	//定义一个未花费的空切片，经过下列操作找出未花费的输出，将其添加到这个切片中，最终返回它
	utxo :=[]*tx.Output{}
	//用于统计全部已花费的输入
	  //map的key是交易的哈希
	  // map的value是输出的索引切片
	spentOuts :=map[string][]int{}
	//思路：要找到为花费的输出，必须要找出所有的输出（遍历输出），根据条件判断就行；
	     //找出所有的输出又必须依赖找到全部的交易;（遍历交易）
	     //而找到全部的交易更加必须依赖找到所有的区块（遍历区块）。
    //下面的三层结构：区块链（区块） -> 区块的交易（每个交易）-> 交易中输出（每个输出）
       //第一层遍历区块（与blockchain.go迭代展示区块中遍历区块语法一致）
    bci :=NewBCIterator(bc)
    for block,err :=bci.Next();err ==nil;block,err =bci.Next(){
    	//第二层遍历交易(for..range..遍历：for下标名，数值名=range 结构体名）
    	for _,tx :=range block.GetTxs() {
    		//第三层遍历所有的输入(输出要判断输出是否会是下一个输入）
    		for _,input :=range tx.Inputs{
    			//记录所有的交易标识（交易hash值）和输出索引
    			   //判断map中的key(hash)是否存在
    			if _,exis :=spentOuts[tx.Hash];!exis{
					// 该交易hash key 不存在，初始化
    				spentOuts[tx.Hash]=[]int{}
				}
				// 该交易key已经存在，追加输出索引即可
				spentOuts[tx.Hash]=append(spentOuts[tx.Hash],input.IndexSrcOutput)
			}

    		//第三层遍历输出
    		for i, output :=range tx.Outputs{
    			//在所有的交易中判断出属于某个人的所有交易
    			//在所有的交易中找到就是未花费的输出，
    			//符合未花费的输出有一个条件：这个输出不会成为下一个输入
    			//这个条件还要依赖找出所有的输入（遍历输入）
    			if output.To ==address && cheUnspent(spentOuts,tx.Hash,i){
    				//满足条件加入最开始初始化的utxo中
    				utxo = append(utxo,output)

				}
			}
		}
	}
	return utxo
}
//将”这个输出不会成为下一个输入“这个条件:花费为false，未花费true i为输出的索引
func cheUnspent(spentOuts map[string][]int,txHash string,i int) bool  {
	indexs,exis :=spentOuts[txHash]
	// 该交易不在全部已花费的输入
	if !exis{
	   	return true
	   }
	// 继续监测索引是否匹配(已花费的是找零还是下一个输入）
	     //找零-一次性只能将整体的UTXO作为另一个交易的输入。类似于纸币的概念，
	         //不能向数字支付一样，给出任意的金额，当纸币的金额超过了需要支付的金额时，就需要找零。
	 for _,index :=range indexs{
	 	//判断本次输出的索引是否是下一个输入的suoyin
	 	if i==index{
	 		return false
		}
	 }
    return true
}
