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
     dec :=gob.NewDecoder(&buffer)
     buffer.Write(data)
     tx :=TX{}
     err :=dec.Decode(&tx);
     if err !=nil {
		 return tx, err
	 }
	return tx,nil
}
