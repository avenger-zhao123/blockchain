package baockchain

import (
	"blockchain/block"
	"blockchain/pow"
	"blockchain/tx"
	"blockchain/wallet"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"time"
)

//一.建立区块链类型结构体
type BlockChain struct {
	lastHash block.Hash //最后一个区块的哈希
	//了解leveldb数据库之后，更新内容
	db *leveldb.DB   //leveldb的数据库连接
	//blocks map[Hash]*Block  //全部区块信息，由区块哈希作为key来检索

	//在建立UTXO结构体时，更新区块链结构体
	UTXOCache *UTXOCache
}
//二.区块链结构体赋值（值可能不为具体的数值）
  //了解leveldb数据库之后，更新内容:给NewBlockchain传递”db leveldb.DB“的参数
func NewBlockchain(db *leveldb.DB)*BlockChain {
	//实例化区块，指在面向对象的编程中，通常把用类创建对象的过程称为实例化
	bc := &BlockChain{
		//blocks : map[Hash]*Block{},  //全部区块信息，由区块哈希作为key来检索
		//了解leveldb数据库之后，更新内容
		db: db,

		///在建立UTXO结构体时，增加UTXOCache
		UTXOCache: NewUTXOCache(db),

	}
	//了解leveldb数据库之后，更新内容：读取存在于数据库的最后区块的哈希值
	//更新原因：读取数据库里最后区块哈希值是为了连接上一个区块链，若不连接，则是新创一个区块链
	//读取key(最后区块的哈希值),详见 "leveldb-test.go"
	data,err :=bc.db.Get([]byte("lastHash"),nil)
	//err为空是指数据库没有lastHash，也就是没有区块，则date赋值于bc.lastHash，若不为空，则直接返回bc
	//若写成err !=nil，则是不执行这两段代码，直接返回bc
	if err ==nil{
		bc.lastHash= block.Hash(data)
	}
	return bc
}
//三.添加创世区块（第一个区块）
//bc *BlockChain是方法，AddGensisBlock是函数名，*BlockChain是返回值
func (bc *BlockChain)AddGensisBlock(address wallet.Address) *BlockChain  {

	//校验是否可以添加创世区块
	if bc.lastHash !="" {  //bc.lastHash已经存在，只是空的字符串，因此不能写成nil
		//数据库存在区块，不需要添加区块，直接返回bc
		return bc
	}
	//数据库不存在区块，需要再添加创世区块
	return bc.AddBlock(address)
}
//四.添加区块到区块链上
//bc *BlockChain是方法，AddBlock是函数名， 提供区块的数据，目前是字符串，*BlockChain是返回值
func (bc *BlockChain) AddBlock(address wallet.Address  ) *BlockChain {
	//构建区块
	b := block.NewBlock(bc.lastHash)
	//在定义交易结构体后,更新增加交易
	//为区块增加交易，任何区块都有coinbase（区块奖励金）交易
	cdtx :=tx.NewCoinbaseTX(address)
	// 将交易加入到区块中
	b.SetTX(cdtx)
//在定义完UTXO结构体后，更新：将新增的交易加入缓存中
	bc.UTXOCache.UpdateUTXO(cdtx)
//在定义完转账交易时，更新：将转账交易加入区块中
    bc.TxsInit()
	//假设每个区块仅仅可以存储3个交易，去掉CoinBase交易，此处从交易缓存中，获取2个交易。
	//实操时，通过交易size在控制的。一个区块为1M，一笔交易大约 250Byte，装满为止。
	for i,l :=0,len(Txs);i<l&&i<2;i++{
		//从交易缓存中，获取交易
		t :=Txs[0]      //让“coinbase（区块奖励金）交易“在交易缓存的最开始
		if len(Txs)>1{   //判断是否有除了“coinbase（区块奖励金）交易”之外的交易
			Txs =Txs[1:] //有就从新定义交易缓存
		}
		// 将交易加入到区块中
		b.SetTX(t)
		// 更新交易对应的 UTXO
		bc.UTXOCache.UpdateUTXO(t)
	}
	bc.TxsSave()
//构造完默克尔树后，更新：添加默克尔树
	// 交易添加完毕，区块已经存储了全部的交易，将全部交易数据提取，构建默克尔树
	merkledata :=[][]byte{}   //定义交易的hash值集合
	for _,t :=range b.GetTxs() {  //遍历所获得交易
		merkledata =append(merkledata,[]byte(t.Hash))  //将获得的交易的hash值加入定义好的集合中
	}
	//调用” 构建默克尔树“函数，将获得交易的hash值传进去
	mt :=block.NewMerleTree(merkledata)
	// 将默克尔树的根节点，放入区块头中
    b.GetHashMerkRoot(mt.Root.Hash)


	//对区块做POW，工作证明
	  //pow对象
	p :=pow.NewPow(b)
	  //开始证明
	nonce,hash := p.Proof()
	  //保险做个判断
	if nonce ==0 || hash == "" {
		log.Fatal("block Hashcash Proof Failed")
	}
	// 为区块设置nonce和hash
	b.SetNonce(nonce).SetHashCurr(hash) //集联调用


	//将区块加入到链的存储结构中
	//bc.blocks[b.hashCurr] =b
	//了解leveldb数据库之后，更新内容:将区块加入到链的存储结构中
	//更新的原因是：当得到区块链实例时，要考虑区块已经存在的情况。 意味着需要去确定最后的区块哈希值。
	//确定最后的区块哈希值方法：是在添加区块时，将最后的区块哈希值储存到数据库中
	if z, err := block.BlockSerialize(*b); err != nil { //建立一个BlockSerialize函数，作用是将Block中的数据转化成byte切片型数据并判断是否有错
		log.Fatal("block can not be serialized.") //出错则返回一段话
		//设置（建立）一个key，key是区块的哈希值，值是上面的byte切片数据
		//key为b_哈希值,加上"b_"是为了标识，通常会在区块hash的key上，增加前缀。
		//设置的时候要考虑是否有错
	} else if err = bc.db.Put([]byte("b_"+b.GetHashCurr()), z, nil);
		err != nil {
		//出错则返回一段话
		log.Fatal("block can not be saved")

	}
	//另一种方法：
	//if z,err :=BlockSerialize(*b); err==nil {
	//	if err =bc.db.Put([]byte("b_"+b.hashCurr),z,nil);err==nil {
	//		bc.lastHash = b.hashCurr
	//		if err = bc.db.Put([]byte("lastHash"), []byte(b.hashCurr), nil); err ==nil{
	//			return bc
	//		}else{
	//			log.Fatal("lasthash can not be saved")
	//			return bc
	//		}
	//	}else {
	//		log.Fatal("block can not be saved")
	//		return bc
	//	}
	//}else{
	//	log.Fatal("block can not be serialized.")
	//	return bc
	//}

	//没有更新的内容：
	//将最后的区块哈希设置为当前区块
	bc.lastHash = b.GetHashCurr()
	// 将最后的区块哈希存储到数据库中
	err := bc.db.Put([]byte("lastHash"), []byte(b.GetHashCurr()), nil)
	if err != nil {
		log.Fatal("lastHas can not be saved")
	}

	return bc
}

//五.通过Hash获取区块
func (bc *BlockChain)GetBlock(hash block.Hash)(*block.Block,error) {
	//从数据库中读取对应的区块
	data,err := bc.db.Get([]byte("b_" + hash),nil)  //key为b_哈希,加上"b_"是为了标识，通常会在区块hash的key上，增加前缀。
	if err !=nil {
		return nil, err
	}
	//反序列化（从数据库读出来是序列化-对应的数据（对象的状态信息）转化成字符串（可以存储或传输的形式），展示是要反序列化的-与序列化相反）
	b,err := block.BlockUnSerialize(data) //在serialize中创建BlockUnSerialize函数，以便调用
	if err !=nil {
		return nil, err
	}
	//函数的返回值是引用型
	return &b,nil
}



//六.迭代展示区块的方法 （方便m命令参数的调用）
func (bc *BlockChain) Iterate() {
	//通过for循环遍历出当前区块
	//for hash :=bc.lastHash;hash !=""; {
	//	//b作为blocks的下标
	//	//b :=bc.blocks[hash]
	//	//得到区块  GetBlock是上面的函数，这块调用
	//	b,err :=bc.GetBlock(hash)
	//	if err!=nil {
	//		log.Fatal(err)
	//		return
	//	}

	//在学完UTXO后更新迭代展示区块的方法
	   // 构建迭代器
	bci :=NewBCIterator(bc)
	  //通过for循环遍历出当前区块
	for b,err :=bci.Next();err ==nil; b,err =bci.Next(){
		//做区块的验证
		   //pow对象 ,调用NewPow，把GetBlock(hash)传进去
		pow := pow.NewPow(b)
		if !pow.Validata(){
			log.Fatalf("Block <%s> is not Valid.",b.GetHashCurr())
			continue  //表示区块非法
		}
        
        	
		//打印区块的Hash值
		fmt.Println("HashCurr:", b.GetHashCurr())
		//打印区块的交易列表
		//fmt.Println("Txs", b.GetTxs())
		//为了更好的展示区块交易的信息（如：hash值，输入，输出）
		fmt.Println("Txs:",b.GetTxsString())
		//打印节点生成时间
		fmt.Println("Time", b.GetTime().Format(time.UnixDate))
		//打印前一个节点的哈希值
		fmt.Println("HashPrev",b.GetHashPrevBlock())
		//hash =b.GetHashPrevBlock()

	}
}



//清空命令
func (bc *BlockChain)Clear() {
	// 数据库中全部区块链的 key 全部删除
	bc.db.Delete([]byte("lastHash"),nil)
	// 迭代删除，全部的 b_ 的key
	//util是leveldb中的，通过b_ key迭代出全部的key
    iter :=bc.db.NewIterator(util.BytesPrefix([]byte("b_")),nil)
    for iter.Next() {
    	bc.db.Delete(iter.Key(),nil)
	}
	//建立好UTXO缓存后，更新：删除UTXO缓存的数据
	iter =bc.db.NewIterator(util.BytesPrefix([]byte("t_")),nil)
	for iter.Next() {
		bc.db.Delete(iter.Key(),nil)
	}
	//j建立好Trabsfer缓存，更新：删除Trabsfer缓存的数据
	bc.db.Delete([]byte("txs"),nil)
	//释放迭代器
	iter.Release()
    //清空bc对象
    bc.lastHash =""
}