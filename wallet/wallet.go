package wallet

import (
	//"crypto/ecdsa"
	//"crypto/elliptic"
	//"crypto/rand"
	"crypto/sha256"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"log"
    "golang.org/x/crypto/ripemd160"
	"github.com/mr-tron/base58"
)

type Address = string
//定义密钥的长度
const keyBitSize = 256
//建立钱包的结构体类型
type Wallet struct {
	//私钥：由椭圆曲线算法构造
	//privatekey *ecdsa.PrivateKey    //定义为引用类型原型是
	//目前公钥作用是生成地址，而公钥可以由私钥生成，因此只需要建立地址就行

	Address Address
//做完助记词测试，更新内容：
	//重新设置私钥
	privatekey *bip32.Key
    //助记词
	mnemonic string
	//重新设置私钥

}
//构造钱包函数
   //更新内容：传参-密码短语
func NewWallet(pass string) *Wallet{
	w :=&Wallet{}
	//生成私钥
	w.Genprivatekey(pass)
	//生成地址
	w.GenAddress()
	return w
}
//利用椭圆曲线算法构造私钥（使用crypto/ecdsa 包完成，go内置包。）
    //私钥是基于随机数生成
func (w *Wallet)Genprivatekey(pass string) *Wallet {
	//利用ecdsa包建立私钥
	//GenerateKey是椭圆曲线算法构造私钥的函数
	   // 参数elliptic.P256() 生成椭圆；参数rand.Reader, 生成随机数
	//privayekey ,err:=ecdsa.GenerateKey(elliptic.P256(),rand.Reader)
	//if err !=nil{
	//	log.Fatal(err)
	//}
	//将生成的私钥赋值给钱包里的私钥
	//w.privatekey =privayekey

//做完助记词测试，更新内容：
	//生成助记词
	   //生成熵
	entropy,err :=bip39.NewEntropy(keyBitSize)
	if err !=nil{
		log.Fatal(err)
	}
	  //基于熵，生成助记词
	mnemonic,err :=bip39.NewMnemonic(entropy)
	if err !=nil{
		log.Fatal(err)
	}
	w.mnemonic =mnemonic
	//基于助记词生成密钥对
	  //基于助记词和短语密码（pass）生成种子
	send :=bip39.NewSeed(mnemonic,pass)
	  //基于种子生成密钥
	privatekey,err :=bip32.NewMasterKey(send)
	if err !=nil {
		log.Fatal(err)
	}
	w.privatekey =privatekey
	return w
}
//生成地址
    //生成公钥（作用：地址基于公钥生成）
    //（椭圆曲线的公钥）方法：，是一个点坐标（椭圆曲线上的点的坐标，平面），通常将X和Y组合在一起形成publicKey。
//func (w *Wallet)genPubkey() []byte {  //定义成字节型切片原因：返回值不需要外部使用，并且做hash运算时传参是字节型切片
//	pubkey :=append(
//		w.privatekey.X.Bytes(),
//		w.privatekey.Y.Bytes()...
//		)
//	return pubkey
//}
 //做完助记词测试，更新内容：
   //生成公钥的hash值
     //将公钥的hash算法单独列出来（相当步骤中第一步）
func HashPubkey(pubkey []byte)[]byte  {
	//步骤一，给公钥做双层hash运算：ripemd160(sha256(pubkey))
	shaHash :=sha256.Sum256(pubkey)
	//调用ripemd160函数
	rpmd := ripemd160.New()
	//将刚做的hash值写入ripemd160函数
	rpmd.Write(shaHash[:])
	//将rpmd的总和赋值给公钥的双层Hash
	pubHash :=rpmd.Sum(nil)
	return pubHash
}
    //生成地址
    //步骤：1. 对 publicKey 做双层的 Hash 运算，得到公钥哈希 pubHash。第一次为 sha256, 第二次为ripeMd160 算法。
         //2. 对公钥哈希再次做双层hash运算，得到 checksum（校验码）。校验码选取前4个字节来使用。双层都为 sha256 运算。
         //3. 将固定的版本号（00）一个字节 与 公钥哈希和校验码组合后，使用Base58编码，形成Address。
func(w *Wallet)GenAddress() *Wallet {
	//将公钥赋值于地址里
	//pubkey :=w.genPubkey()
	////步骤一，给公钥做双层hash运算：ripemd160(sha256(pubkey))
	//shaHash :=sha256.Sum256(pubkey)
	////调用ripemd160函数
	//rpmd := ripemd160.New()
	////将刚做的hash值写入ripemd160函数
	//rpmd.Write(shaHash[:])
	////将rpmd的总和赋值给公钥的双层Hash
    //pubHash :=rpmd.Sum(nil)
//做完助记词测试，更新内容：
	//将公钥建立出来
	pubkey :=w.privatekey.PublicKey().String()
	//调用公钥hash算法的函数，将公钥传进去
	hashPubkey :=HashPubkey([]byte(pubkey))
    //步骤二：得到校验码
	//对公钥哈希再次做双层hash运算(都为sha256）
    //为更新时：a :=sha256.Sum256(pubHash)
    a :=sha256.Sum256(hashPubkey)
    checksum :=sha256.Sum256(a[:])  //a[:]表示数组型转化为切片
    //步骤三：形成地址
       //将版本号、公钥哈希值、校验码组合
    //未更新时：c :=append([]byte{0},pubHash...)
	c :=append([]byte{0},hashPubkey...)
    data :=append(c,checksum[:2]...)   //校验码选取前4个字节来使用(不包括第四个字节）
      //使用Base58编码，形成Address
    w.Address =base58.Encode(data)
    //返回钱包
     return w
}
//获取助记词
func (w *Wallet)GetMnemonic() string {
	return w.mnemonic
	
}