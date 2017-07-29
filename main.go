package main

import (
	"fmt"
	"time"
	"crypto/sha256"
	"encoding/binary"
)

type Block struct {
	index uint64
	timestamp time.Time
	data string
	previousHash []byte
	hash []byte
}

func NewBlock(
	index uint64,
	timestamp time.Time,
	data string,
	previousHash []byte) *Block {

	block := new(Block)
	block.index = index
	block.timestamp = timestamp
	block.data = data
	block.previousHash = previousHash

	h := sha256.New()

	// Index
	indexBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(indexBytes, index)
	h.Write(indexBytes)

	// Timestamp
	timestampBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestampBytes, uint64(timestamp.UnixNano()))
	h.Write(timestampBytes)

	// Data
	h.Write([]byte(data))

	// Previous hash
	h.Write(previousHash)

	block.hash = h.Sum(nil)

	return block
}

func main() {
	genesisBlock := NewBlock(0, time.Now(), "Genesis Block", make([]byte, 0))

	fmt.Printf("%x\n", genesisBlock.hash)
}
