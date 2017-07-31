package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
	"net/http"
)

type Block struct {
	index        uint64
	timestamp    time.Time
	data         string
	previousHash []byte
	hash         []byte
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

// NextBlock returns next block
func NextBlock(lastBlock *Block) *Block {
	return NewBlock(
		lastBlock.index+1,
		time.Now(),
		fmt.Sprintf("Block number: %d", lastBlock.index+1),
		lastBlock.hash)
}

func main() {
	genesisBlock := NewBlock(0, time.Now(), "Genesis Block", make([]byte, 0))

	fmt.Printf("%x\n", genesisBlock.hash)

	blockchain := []*Block{genesisBlock}
	previousBlock := genesisBlock

	for i := 0; i < 20; i++ {
		newBlock := NextBlock(previousBlock)
		blockchain = append(blockchain, newBlock)
		previousBlock = newBlock

		fmt.Printf("Block hash %x\n", newBlock.hash)
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Input: %s!", r.URL.Path[1:])
}
