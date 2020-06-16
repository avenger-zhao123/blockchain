package tx


//定义输入的结构体类型
type Input struct {
	// 输入来源交易(上一笔交易：比如B的钱是A给的，B的交易为当前交易，那A的交易为上一笔交易)的hash
	HahSrcTX string
	// 输入来源交易输出的索引
	//来源交易输出的索引(B的钱是A给的，A可能不止给B发钱，需要确定B从A得到的是哪一部份钱）
	IndexSrcOutput string
}
