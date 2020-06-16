package block

import (
	"blockchain/tx"
	"bytes"
	"encoding/gob"
	"time"
)

// 由于区块的字段都是 unexported field 非导出字段
//使用中间的数据结构作为桥梁，完成序列化。（就是将区块和区块头展示在同一结构体上）
type BlockData struct {
	//为何要在序列化中再建立一个结构体，而不是将区块和区块头里的字段开头大写？
	//原因是若将区块和区块头里的字段开头大写，任何人都可以改字段，而在序列化里只能查看，不能改
	Version        int
	HashPrevBlock  Hash
	HashMerkleRoot Hash
	Time           time.Time
	Bits           int
	Nonce          int
	Txs            []*tx.TX
	TxCounter      int
	HashCurr       Hash
}
//区块序列化
func BlockSerialize(b Block)([]byte,error) {
	// 将区块数据赋值到 BlockData
	e := BlockData{
		Version: b.header.version,
		HashPrevBlock: b.header.hashPrevBlock,
		HashMerkleRoot: b.header.hashMerkRoot,
		Time: b.header.time,
		Bits: b.header.bits,
		Nonce: b.header.nonce,
		Txs: b.txs,
		TxCounter: b.txCouner,
		HashCurr: b.hashCurr,
	}
	//执行 gob 序列化即可(gob编码）详见"gob-test.go"
     //创建容器
     buffer :=bytes.Buffer{}
     //创建编码器
     enc :=gob.NewEncoder(&buffer)
     //编码数据（序列化数据）
     err :=enc.Encode(e)
     if err !=nil{
     	return nil, err
	 }
	 //编码成功
	 return buffer.Bytes(),nil

}
//区块反序列化（解码）
func BlockUnSerialize(data []byte)(Block,error) {
	//与解码的步骤异曲同工，详见"gob-test.go"
	//建立容器
	unbuffer :=bytes.Buffer{}
	//将之前的编码数据，放入缓存中
	unbuffer.Write(data)
	//建立解码器
	dec :=gob.NewDecoder(&unbuffer)
	//解码时需要提供解码的数据类型
	f := BlockData{}
	//解码数据(反序列化）
	err:=dec.Decode(&f)
	if err  !=nil{
		return Block{}, err
	}
	// 反序列化成功
	return Block{
		header: BlockHeader{
			version: f.Version,
			hashPrevBlock: f.HashPrevBlock,
			hashMerkRoot: f.HashMerkleRoot,
			time: f.Time,
			bits: f.Bits,
			nonce: f.Nonce,
		},
		txs: f.Txs,
		txCouner: f.TxCounter,
		hashCurr: f.HashCurr,
	}, nil
}
