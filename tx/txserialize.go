package tx

import (
	"bytes"
	"encoding/gob"
)

//序列化TX(与区块的序列化函数基本一致）
func SerializeTX(tx TX) ([]byte,error)  {
	//执行 tx 序列化即可(tx编码）详见"gob-test.go"
	//创建容器
     buffer :=bytes.Buffer{}
	//创建编码器
     enc :=gob.NewEncoder(&buffer)
	//编码数据（序列化数据）
     err :=enc.Encode(tx)
     if err !=nil {
     	return nil, err
	  }
	//编码成功
	  return buffer.Bytes(),nil
}
//反序列化（反穿行化，解码）(与区块的反序列化函数基本一致）
func UnserializaTX(data []byte)(TX,error)  {
     buffer :=bytes.Buffer{}
	//将之前的编码数据，放入缓存中
     dec :=gob.NewDecoder(&buffer)
	//建立解码器
     buffer.Write(data)
	//解码时需要提供解码的数据类型
     tx :=TX{}
	//解码数据(反序列化）
     err :=dec.Decode(&tx)
     if err !=nil {
		 return tx, err
	 }
	return tx,nil
}
