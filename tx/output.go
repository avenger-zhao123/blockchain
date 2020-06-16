package tx

import "blockchain/wallet"

//定义输出的结构体类型
type Output struct {
	//输出的金额
	Value int
	// 目标用户地址
	To wallet.Address
}
