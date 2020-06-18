package baockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

//建立UTXO序列化函数
func SerializeUTXOSet (us UTXOSet) ([]byte,error) {
	//建立一个容器
	buffrer :=bytes.Buffer{}
	//进行编码
	enc :=gob.NewEncoder(&buffrer)
	//判断编码是否有错
	err :=enc.Encode(&us)
	if err !=nil {
		log.Fatal(err)
	}
	//返回序列化的数和空错误
	return buffrer.Bytes(),nil
}
//建立UTXO反序列化函数
func UnSerializeUTXOSet(data []byte)  (UTXOSet,error) {
	//建立容器
	buffer :=bytes.Buffer{}
	//将之前的编码数据，放入缓存中
	buffer.Write(data)
	//建立解码器
	dec :=gob.NewDecoder(&buffer)
	//解码时需要提供解码的数据类型
	us :=UTXOSet{}
	//解码数据(反序列化）
	 err :=dec.Decode(&us)
	 if err !=nil{
	 	log.Fatal(err)
	 }
	return us,err
	
}