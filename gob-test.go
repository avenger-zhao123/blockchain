package main

import (
	"bytes"
	"encoding/gob"
	 "fmt"
	"log"
)

type Block struct {
	CurrHash string
	Txs string
}
func main()  {
	b :=Block{
		CurrHash:"228422940",
		Txs:"first transaction",
	}
     //gob编码
      //编码前需要一个编码器，因此要先创建一个编码器
       //编码器则需要一个能写入内容的容器
	var encodercontainer bytes.Buffer  //容器为bytes.Buffer（byte 型的缓存（就是变量））提供的缓存，应该具备可写功能
	  //创建一个编码器（Encoder)
	  //编码器将来将编码的结果写入缓存中
	enc :=gob.NewEncoder(&encodercontainer)
	  //编码数据，编码的结果写入编码器的缓存中
	 err :=enc.Encode(b)
	 if err!=nil{
	 	log.Fatal(err)

	}
	fmt.Println(encodercontainer.Bytes())
	//将编码的数据定义为一个变量，以便解码
	a :=encodercontainer.Bytes()


	//gob解码
	  //解码前同样需要一个解码器
	    //解码器同样也需要一个能够读取的容器
	    // 提供的缓存，应该具备可读功能
	var  decodercontainer bytes.Buffer
	//将之前的编码数据，放入缓存中
	decodercontainer.Write(a)
	//创建一个解码器（Decoder)
	//解码器将来将编码的数据从缓存中解码出来
	dec :=gob.NewDecoder(&decodercontainer)
	//解码时需要提供解码的数据类型
	c :=Block{}
	//解码数据
	err =dec.Decode(&c)
	if err !=nil {
		log.Fatal(err)
	}

	fmt.Println(c)


}
