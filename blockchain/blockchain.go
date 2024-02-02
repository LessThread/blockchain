package blockchain

import (
	"encoding/hex"
	"fmt"
	"goblockchain/constcoe"
	"goblockchain/transaction"
	"goblockchain/utils"
	"runtime"

	"github.com/dgraph-io/badger"
)

// 链定义
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

// 初始化
func InitBlockChain(address []byte) *BlockChain {
	var LastHash []byte
	if utils.FileExists(constcoe.BCFile) {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(constcoe.BCPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	utils.Handle(err)

	err = 
}

// add block to blockchain
func (bc *BlockChain) AddBlock(txs []*transaction.Transaction) {
	newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, txs)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// 在整个链上寻找前置交易信息(UTXO)
func (bc *BlockChain) FindUnspentTransactions(address []byte) []transaction.Transaction {
	var unSpentTxs []transaction.Transaction
	spentTxs := make(map[string][]int) // can't use type []byte as key value
	// 从新区块遍历整个区块链
	for idx := len(bc.Blocks) - 1; idx >= 0; idx-- {
		block := bc.Blocks[idx]
		// 对于每个区块，遍历其交易
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// 对于每个交易，遍历交易信息中的OutPut
		IterOutputs:
			for outIdx, out := range tx.Outputs {

				//检查是否被使用过
				if spentTxs[txID] != nil {
					for _, spentOut := range spentTxs[txID] {
						if spentOut == outIdx {
							continue IterOutputs
						}
					}
				}

				// 统计支付方在之前交易中的入账
				if out.ToAddressRight(address) {
					unSpentTxs = append(unSpentTxs, *tx)
				}
			}

			// 检查是否是创世交易，否则遍历Inuput
			if !tx.IsBase() {
				for _, in := range tx.Inputs {

					// 索引支付方地址在之前交易过程中的已经用掉的（统计出账）
					if in.FromAddressRight(address) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTxs[inTxID] = append(spentTxs[inTxID], in.OutIdx)
					}
				}
			}
		}
	}
	return unSpentTxs
}

// 计算余额
func (bc *BlockChain) FindUTXOs(address []byte) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) {
				accumulated += out.Value
				unspentOuts[txID] = outIdx
				continue Work // one transaction can only have one output referred to adderss
			}
		}
	}
	return accumulated, unspentOuts
}

// 找到足够支付的UTXO
func (bc *BlockChain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = outIdx
				if accumulated >= amount {
					break Work
				}
				continue Work // one transaction can only have one output referred to adderss
			}
		}
	}
	return accumulated, unspentOuts
}

// 创建交易项
func (bc *BlockChain) CreateTransaction(from, to []byte, amount int) (*transaction.Transaction, bool) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		fmt.Println("Not enough coins!")
		return &transaction.Transaction{}, false
	}

	// 将预选UTXO加入本次Transaction的Input
	for txid, outidx := range validOutputs {
		txID, err := hex.DecodeString(txid)
		utils.Handle(err)
		input := transaction.TxInput{TxID: txID, OutIdx: outidx, FromAddress: from}
		inputs = append(inputs, input)
	}

	outputs = append(outputs, transaction.TxOutput{Value: amount, ToAddress: to})

	// 找零
	if acc > amount {
		outputs = append(outputs, transaction.TxOutput{Value: acc - amount, ToAddress: from})
	}
	tx := transaction.Transaction{ID: nil, Inputs: inputs, Outputs: outputs}
	tx.SetID()

	return &tx, true
}

// 封装用函数，将交易打包进入内存池
func (bc *BlockChain) Mine(txs []*transaction.Transaction) {
	bc.AddBlock(txs)
}
