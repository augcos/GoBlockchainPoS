package blockchainPoS

import (
	"time"
	"crypto/sha256"
	"encoding/hex"
)

// var Blockchain contains the multiple Blocks
var Blockchain []Block
type Block struct {
	BlockNumber int
	BlockTime 	string
	Data		string
	Hash		string
	PrevHash	string
	Validator	string
}

// CalculateHash() returns the sha256 of a given string
func CalculateHash(s string) string{
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	hash := hex.EncodeToString(hashed)
	return hash
}

// CalculateBlockHash() returns the hash of a given block
func CalculateBlockHash(block Block) string{
	record := string(block.BlockNumber) + block.PrevHash
	return CalculateHash(record)
}

// GenerateBlock() creates a new block given the previous block, a string and the validator's address
func GenerateBlock(prevBlock Block, data string, address string) (Block, error) {
	var nextBlock Block
	t := time.Now()

	nextBlock.BlockNumber = prevBlock.BlockNumber + 1
	nextBlock.BlockTime = t.String()
	nextBlock.Data = data
	nextBlock.PrevHash = prevBlock.Hash
	nextBlock.Validator = address
	nextBlock.Hash = CalculateBlockHash(nextBlock)

	return nextBlock, nil
}

// IsBlockValid() checks if a block's BlockNumber, PrevHash and Hash are correct
func IsBlockValid(prevBlock, nextBlock Block) bool {
	if prevBlock.BlockNumber+1 != nextBlock.BlockNumber {
		return false
	}
	if prevBlock.Hash != nextBlock.PrevHash {
		return false
	}
	if CalculateBlockHash(nextBlock) != nextBlock.Hash {
		return false
	}
	return true
}

// IsBlockchainValid() checks if all blocks in a blockchain are valid
func IsBlockchainValid(testedChain []Block) bool {
	for i:=1; i<(len(testedChain)-1); i++ {
		if !IsBlockValid(testedChain[i-1],testedChain[i]) {
			return false
		}
	}
	return true
}

// ReplaceChain() replaces the blockchain if it is given a longer alternative
func ReplaceChain(nextBlockchain []Block) {
	if len(nextBlockchain) > len(Blockchain) {
		Blockchain = nextBlockchain
	}
}