package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"goblockchain/transaction"
	"goblockchain/utils"
	"time"
)

type Block struct {
	Timestamp    int64
	Hash         []byte
	PrevHash     []byte
	Target       []byte
	Nonce        int64
	Transactions []*transaction.Transaction
}

// 创造区块
func CreateBlock(prevhash []byte, txs []*transaction.Transaction) *Block {
	block := Block{time.Now().Unix(), []byte{}, prevhash, []byte{}, 0, txs}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

// 创世区块
func GenesisBlock(address []byte) *Block {
	tx := transaction.BaseTx([]byte(address))
	genesis := CreateBlock([]byte("block_god"), []*transaction.Transaction{tx})
	genesis.SetHash()
	return genesis
}

// 返回交易hash集合
func (b *Block) BackTrasactionSummary() []byte {
	txIDs := make([][]byte, 0)
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.ID)
	}
	summary := bytes.Join(txIDs, []byte{})
	return summary
}

// 计算区块哈希
func (b *Block) SetHash() {
	information := bytes.Join([][]byte{utils.ToHexInt(b.Timestamp), b.PrevHash, b.Target, utils.ToHexInt(b.Nonce), b.BackTrasactionSummary()}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

// 区块字节化编码
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	utils.Handle((err))
	return res.Bytes()
}

// 区块解码
func DeSerializeBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader((data)))
	err := decoder.Decode(&block)
	utils.Handle(err)
	return &block
}
