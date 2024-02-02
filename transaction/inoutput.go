package transaction

import "bytes"

type TxOutput struct {
	Value     int    //转出的资产价值
	ToAddress []byte //接收者地址
}

type TxInput struct {
	TxID        []byte //指明支持本次交易的前置交易信息
	OutIdx      int    //具体指明是前置交易信息中的第几个Output
	FromAddress []byte //资产转出者的地址
}

func (in *TxInput) FromAddressRight(address []byte) bool {
	return bytes.Equal(in.FromAddress, address)
}

func (out *TxOutput) ToAddressRight(address []byte) bool {
	return bytes.Equal(out.ToAddress, address)
}
