package baockchain

import (
	"blockchain/tx"
	"blockchain/wallet"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
)
//在考虑性能的情况下，下面的遍历区块就成本很高，若区块有成百上千的话就会高到不现实，
//解决的办法是：将UTXO缓存起来（就相当于给UTXO加上索引），一旦UTXO发生改变，就更新UTXO缓存
//定义一个单独的UTXO结构体类型
type UTXO struct {
	//定义输出，可以引用tx的output
	Output *tx.Output
	//下面的相当于定义了一个输入
	//定义所属交易的Hash值
	HashStrTX string
	//位于交易输出列表的索引
	IndexStrOutput int
}
//定义一个单个交易的UTXO集合
type UTXOSet =[]*UTXO

//建立UTXO 缓存操作对象（对象相当于leveldb数据库）结构体
type UTXOCache struct {
	db *leveldb.DB
}
//UTXO 缓存操作对象赋上leveldb数据库
func NewUTXOCache(db *leveldb.DB) *UTXOCache  {
	return &UTXOCache{
		db :db,
	}
}
//建立UTXO的缓存
func (uc *UTXOCache)UpdateUTXO(tx *tx.TX) *UTXOCache {
	//在 UTXO 缓存中，加入新增的 UTXO
	//步骤为:
	//1.遍历tx交易全部的输出，每个输出作为一个新的UTXO缓存来使用即可
	//(1)定义UTXO的集合，以便之后将UTX放进去
	us :=UTXOSet{}
	//(2)遍历交易
	for i,output :=range tx.Outputs{
		//将UTXO放入集合中
		us =append(us,&UTXO{
			Output: output,
			HashStrTX: tx.Hash,
			IndexStrOutput: i,
		})
	}
	//2.更新对应的缓存
	  //(1)调用序列化函数，将us传进去
	ser,err :=SerializeUTXOSet(us)
	    //判断序列化是否有错
	if err !=nil{
		log.Fatal("UTXOSet Serialize failed")
	}
	  //(2)设置（创建，更新）key（leveldb是key-value型数据库）
	      //key-""t_"+tx.Hash";数值-“ser”，选项-“nil"
	uc.db.Put([]byte("t_"+tx.Hash),ser,nil)
   return uc
}

//构造找到UTXO的函数（功能与”找到哪些属于address的未被花费的输出Unspend“的函数一致）
func(uc *UTXOCache) FindUTXO(address wallet.Address)[]*UTXO {
	//步骤：
	//1.定义一个未花费的空切片，经过下列操作找出未花费的输出，将其添加到这个切片中，最终返回它
    utxo:=UTXOSet{}
	//2.遍历utxo缓存，得到utxo即可(详见blockchain中的Clear()函数)
    iter :=uc.db.NewIterator(util.BytesPrefix([]byte("t_")),nil)
    for iter.Next() {
    	//(1)获取所有的UTXO
    	value,err :=uc.db.Get(iter.Key(),nil)
    	if err !=nil{
    		continue
		}
		//(2)反序列化
		//原因:由于value是序列化的UTXOSet,要想获取Value的数据，先反序列化
		us,err :=UnSerializeUTXOSet(value)
		if err !=nil{
			log.Print(err)
			continue
		}
		//(3)找us中属于address的加入整体的UTXO集合中
		for _,u :=range us{
			if u.Output.To == address {
				utxo =append(utxo,u)
			}
		}
	}
    iter.Release()
    return utxo
}



// 找到哪些属于address的未被花费的输出Unspend
   //UTXO：没有作为其他交易输入的输出，称为为花费的交易输出。
//func (bc *BlockChain)FindUTXO(address wallet.Address) []*tx.Output  {
//	//定义一个未花费的空切片，经过下列操作找出未花费的输出，将其添加到这个切片中，最终返回它
//	utxo :=[]*tx.Output{}
//	//用于统计全部已花费的输入
//	  //map的key是交易的哈希
//	  // map的value是输出的索引切片
//	spentOuts :=map[string][]int{}
//	//思路：要找到为花费的输出，必须要找出所有的输出（遍历输出），根据条件判断就行；
//	     //找出所有的输出又必须依赖找到全部的交易;（遍历交易）
//	     //而找到全部的交易更加必须依赖找到所有的区块（遍历区块）。
//    //下面的三层结构：区块链（区块） -> 区块的交易（每个交易）-> 交易中输出（每个输出）
//       //第一层遍历区块（与blockchain.go迭代展示区块中遍历区块语法一致）
//    bci :=NewBCIterator(bc)
//    for block,err :=bci.Next();err ==nil;block,err =bci.Next(){
//    	//第二层遍历交易(for..range..遍历：for下标名，数值名=range 结构体名）
//    	for _,tx :=range block.GetTxs() {
//    		//第三层遍历所有的输入(输出要判断输出是否会是下一个输入）
//    		for _,input :=range tx.Inputs{
//    			//记录所有的交易标识（交易hash值）和输出索引
//    			   //判断map中的key(hash)是否存在
//    			if _,exis :=spentOuts[tx.Hash];!exis{
//					// 该交易hash key 不存在，初始化
//    				spentOuts[tx.Hash]=[]int{}
//				}
//				// 该交易key已经存在，追加输出索引即可
//				spentOuts[tx.Hash]=append(spentOuts[tx.Hash],input.IndexSrcOutput)
//			}
//
//    		//第三层遍历输出
//    		for i, output :=range tx.Outputs{
//    			//在所有的交易中判断出属于某个人的所有交易
//    			//在所有的交易中找到就是未花费的输出，
//    			//符合未花费的输出有一个条件：这个输出不会成为下一个输入
//    			//这个条件还要依赖找出所有的输入（遍历输入）
//    			if output.To ==address && checkUnspent(spentOuts,tx.Hash,i){
//    				//满足条件加入最开始初始化的utxo中
//    				utxo = append(utxo,output)
//
//				}
//			}
//		}
//	}
//	return utxo
//}
//将”这个输出不会成为下一个输入“这个条件:花费为false，未花费true i为输出的索引
func checkUnspent(spentOuts map[string][]int,txHash string,i int) bool  {
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






