package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
	"net/http"
	"encoding/json"
)

type Block struct {
	index        uint64
	timestamp    time.Time
	data         string
	previousHash []byte
	hash         []byte
}

// NewBlock creates new block
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

	http.HandleFunc("/transaction", transactionHandler)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func transactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		readBuffer := make([]byte, 1024, 1024*128)
		var f interface{}
		readerCloser := r.Body
		count, err := readerCloser.Read(readBuffer)

		fmt.Printf("Body size: %d\n", count)
		fmt.Printf(string(readBuffer), "\n")

		if err != nil && err.Error() != "EOF" {
			fmt.Println("Read error: ", err.Error())
			readerCloser.Close()
			return
		}

		err = json.Unmarshal(readBuffer, &f)

		fmt.Println(f)

		if err != nil {
			fmt.Println("Parse error: ", err.Error())
			readerCloser.Close()
			return
		}

		readerCloser.Close()
	} else {
		fmt.Fprint(w, `{"success": false, "error": "Only POST method is supported"}`)
	}

}
