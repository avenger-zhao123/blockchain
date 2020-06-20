package baockchain

import (
	"blockchain/tx"
	"blockchain/wallet"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
)
//交易缓存
var Txs =[]*tx.TX{}

//建立转账交易函数
func (bc *BlockChain)Transfer(from, to wallet.Address,value int ) error  {
	//1.构建交易数据
	    //(1).查询当前用户的全部UTXO
	utxos := bc.UTXOCache.FindUTXO(from)
	    //(2).凑够转账金额:调用“找到可使用的UTXO”，完成此操作
	amount,spendableUtxos :=FindSpendableUTXO(utxos,value)
	    //(3).判断凑的金额是否够
	if amount <value{
		   //不够直接输出
		return errors.New(fmt.Sprintf("Balance of %s is not enough.",from))
	}
	   //金额够，建立交易的输入和输出
	    //（4）建立输入
	inConter :=len(spendableUtxos)   //测出可使用的UTXO集合长度
	inputs:=make([]*tx.Input,inConter)  //规定输入的长度
	for i,u :=range spendableUtxos{    //遍历可使用的UTXO集合
		inputs[i] =&tx.Input{         //将可使用的UTXO集合中的输入加入到建立的输入中
			HashSrcTX: u.HashStrTX,
			IndexSrcOutput: u.IndexStrOutput,
		}
	}
	    //（5）建立输出
	outputs :=[]*tx.Output{   //输出就是交易中输出
		&tx.Output{           //将参数赋值给对应的字段
			Value: value,
			To: to,
		},
	}
	     //当凑得金额大于所要的，需要找零
	if amount >value {
		outputs =append(outputs,&tx.Output{  //找零操作
			Value: amount-value,             //金额是凑得金额减去所需的金额
			To: from,                        //对象变成付金额的地址
		})
	}
	   //（6）构建好了交易
	t :=tx.NewTransferTX(inputs,outputs)
	//2.储存交易(存储在leveldb上）
	  //（1）获取当前的txs
	bc.TxsInit()
	  //（2）追加
	Txs = append(Txs,t)
	  //（3）Txs的持久化
	bc.TxsSave()

    return nil
}
//构造“找到可使用的UTXO”的函数
func FindSpendableUTXO(utxos []*UTXO, value int)(int,[]*UTXO)  {
  //1.初始化
	//（1）初始化一个凑得金额
	amount:=0
	//（2）初始化一个凑够的UTXO集合
	sautxos :=[]*UTXO{}
  //2.凑金额
	//遍历所有的UTXO
	for _,u :=range utxos{
		//找到每笔交易输出中的金额并把它们加入金额中
		amount+=u.Output.Value
		//将找到的utxo放到集合中
		sautxos=append(sautxos,u)
		//当凑的金额大于所需要的金额，输出
		if amount >=value {
			break
		}

	}
    return amount,sautxos
}
//Txc持久化
func (bc *BlockChain)TxsSave() error {
	//序列化
	buffer :=bytes.Buffer{}
	enc :=gob.NewEncoder(&buffer)
	err :=enc.Encode(Txs)
	if err !=nil {
		return err
	}
	//存储
	err =bc.db.Put([]byte("txs"),buffer.Bytes(),nil)
	if err !=nil{
		return err
	}
    return nil
}
//Txs初始化
func (bc *BlockChain)TxsInit() error  {
	//获取
	 date,err :=bc.db.Get([]byte("txs"),nil)
	if err !=nil{
		return err
	}
	//反序列化
	buffer :=bytes.Buffer{}
	buffer.Write(date)
	dec :=gob.NewDecoder(&buffer)
	err =dec.Decode(&Txs)
	if err !=nil{
		return err
	}

	return nil
}