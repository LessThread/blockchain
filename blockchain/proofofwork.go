package blockchain

import (
	"bytes"
	"crypto/sha256"
	"goblockchain/constcoe"
	"goblockchain/utils"
	"math"
	"math/big"
)

// 获取区块难度
func (b *Block) GetTarget() []byte {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-constcoe.Difficulty))
	return target.Bytes()
}

// 构造含Nonce的区块
func (b *Block) GetBase4Nonce(nonce int64) []byte {
	data := bytes.Join([][]byte{
		utils.ToHexInt(b.Timestamp),
		b.PrevHash,
		utils.ToHexInt(int64(nonce)),
		b.Target,
		b.BackTrasactionSummary(),
	},
		[]byte{},
	)
	return data
}

// 计算Nonce的合法性
func (b *Block) FindNonce() int64 {
	var intHash big.Int
	var intTarget big.Int
	var hash [32]byte
	var nonce int64
	nonce = 0
	intTarget.SetBytes(b.Target)

	for nonce < math.MaxInt64 {
		data := b.GetBase4Nonce(nonce)
		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(&intTarget) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce
}

// 验证hash合法性
func (b *Block) ValidatePoW() bool {
	var intHash big.Int
	var intTarget big.Int
	var hash [32]byte
	intTarget.SetBytes(b.Target)
	data := b.GetBase4Nonce(b.Nonce)
	hash = sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	if intHash.Cmp(&intTarget) == -1 {
		return true
	}
	return false
}
