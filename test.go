package main

import (
	"fmt"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip32"
)

func main(){
	//区块测试
//	c :=baockchain.NewBlock("","Gensis Block.")
//	fmt.Println(c)
//	//数据库连接
//	d :="data"
//	db,err :=leveldb.OpenFile(d,nil)
//	if err !=nil{
//		log.Fatal(err)
//	}
//	//释放数据库连接
//	defer db.Close()
//
//    //区块链测试
//	a :=baockchain.NewBlockchain(db)
//	a.AddGensisBlock()
//	a.AddBlock("First Block").AddBlock("Second Block")
//	//fmt.Println(a)
//	a.Iterate()
//


//挖矿的简易测试
	//定义难易程度
//	bits :=8
//	//定义一个目标（即比它小的数达到工作量证明）
//	target :=big.NewInt(1) //初始为00000 ... 0001
//	// 采用左移位的方案，构建目标比较多
//	//00000001 LSH 1 = 00000010
//	//00000001 LSH 2 = 00000100
//	target.Lsh(target,uint(256-bits+1))   //256-bits+1保证目标和难易程度保持一致
//	//打印出目标数
//	fmt.Println(target.String())
//	fmt.Println("-------Minting----------")
//	//定义一个服务字符串（除随机数之外的区块信息）
//    serviceStr :="block data"
//    //定义一个类型为big.int类型的数据
//      //big.int表示任意长度的整型
//      //使用big.int原因：区块链中，hash 算法为 sha256 算法，意味着hash值的长度为256bits
//
////hash证明
//    var hashInt big.Int
//    nonce :=0  //定义一个计数器
//    //for循环：初始化为：hashInt，条件为循环体中的if语句，条件变化是nonce
//    for{
//        //在哈希之前要将服务器字符串(serviceStr)和随机字符串（nonce计数器）连接
//    	  //trconv.Itoa()函数的参数是一个整型数字，它可以将数字转换成对应的字符串类型的数字
//    	g :=serviceStr +strconv.Itoa(nonce)
//    	//计算g的哈希值，由于sha256.Sum256传参为字节型切片，要将g转化为此类型
//    	hash :=sha256.Sum256([]byte(g))
//    	//将计算好的hash值传到hashInt
//    	  //hash[:]表示数组型转化为切片
//    	   //转化原因：sha256.Sum256函数返回一个数组类型，而hashInt.SetBytes函数的传参是切片类型
//    	hashInt.SetBytes(hash[:])
//    	fmt.Println(hashInt.String(),nonce)
//    	//循环条件，cmp为比较，是将hashInt和target比较
//    	   //若结果为-1，则hashInt < target；若结果为1，则ashInt > target
//    	if hashInt.Cmp(target) ==-1 {
//    		fmt.Printf("挖矿成功")
//			return
//		}
//		//条件变化
//        nonce++
//	}

//助记词测试
	  //目前会生成秘钥使用的助记词。由 BIP39 提供的助记词；助记词的与秘钥一一对应的关系。
	  //秘钥的二进制数据每11bits一组。每组与一个助记词相互对应。若使用的256bits的秘钥，则每组助记词应该是24个（256/11=23.2··==24）
	    //先生成熵（混乱程度），提供特定的比特数
	entropy,_ :=bip39.NewEntropy(256)
	   //基于熵，生成助记词
	mnemonic,_ :=bip39.NewMnemonic(entropy)
	fmt.Println(mnemonic)
    //基于助记词生成密钥对
      //基于助记词和短语密码（用户需要指定）生成种子
   send :=bip39.NewSeed(mnemonic,"Secret Passphrase")
      //基于种子生成密钥
   masterkey,_ :=bip32.NewMasterKey(send)
      //构建公钥（私钥生成）
   publickey :=masterkey.PublicKey()
   fmt.Println("PrivateKey: ", masterkey.String())
   fmt.Println("publickey: ",publickey.String())
}
