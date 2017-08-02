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
	proofOfWork  int
	hash         []byte
}

var blockchain = make([]*Block, 0, 1024*128)
var transactions = make([]interface{}, 0, 1024*128)

// NewBlock creates new block
func NewBlock(
	index uint64,
	timestamp time.Time,
	data string,
	proofOfWork int,
	previousHash []byte) *Block {

	block := new(Block)
	block.index = index
	block.timestamp = timestamp
	block.data = data
	block.proofOfWork = proofOfWork

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
	genesisBlock := NewBlock(0, time.Now(), `{}`, 2, make([]byte, 0))

	fmt.Printf("%x\n", genesisBlock.hash)

	blockchain = append(blockchain, genesisBlock)

	http.HandleFunc("/transaction", transactionHandler)
	http.HandleFunc("/mine", mineHandler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func transactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprint(w, `{"success": false, "error": "Only POST method is supported"}`)
		return;
	}

	readBuffer := make([]byte, 1024, 1024*128)

	readerCloser := r.Body
	count, err := readerCloser.Read(readBuffer)

	body := readBuffer[:count]

	fmt.Printf("Body size: %d\n", count)
	fmt.Println(string(body))

	if err != nil && err.Error() != "EOF" {
		fmt.Println("Read error: ", err.Error())
		readerCloser.Close()
		return
	}

	var transactionMap interface{}
	err = json.Unmarshal(body, &transactionMap)

	fmt.Println(transactionMap)

	if err != nil {
		fmt.Println("Parse error: ", err.Error())
		readerCloser.Close()
		return
	}

	readerCloser.Close()

	transactions = append(transactions, transactionMap)

	fmt.Println("Transactions")
	fmt.Println(transactions)

	fmt.Fprint(w, `{"success": true}`)
}

func proofOfWork(lastProof int) int {
	fmt.Printf("lastProof: %d\n", lastProof)
	inc := lastProof + 1

	fmt.Printf("inc: %d\n", inc)
	for ; !((inc % 9 == 0) && (inc % lastProof == 0));  inc += 1 {}

	fmt.Println("New proof of work found: %d", inc)

	return inc
}

func mineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fmt.Fprint(w, `{"success": false, "error": "Only GET method is supported"}`)
		return;
	}

	lastBlock := blockchain[len(blockchain) - 1]

	proof := proofOfWork(lastBlock.proofOfWork)

	newTransaction:= `{"from": "network", "to": miner_address, "amount": 1}`
	transactions = append(transactions, newTransaction)

	newBlockData :=  fmt.Sprintf(`{"transactions": "%v" }`, transactions)

	newBlock := NewBlock(lastBlock.index + 1, time.Now(), newBlockData, proof, lastBlock.hash)
	blockchain = append(blockchain, newBlock)

	// Empty transaction list
	transactions = make([]interface{}, 0, 1024*128)

	// Send new block to client
	fmt.Fprint(w, newBlock.reducedJSON())
}

func (block *Block) reducedJSON() string {
	return fmt.Sprintf(`{` +
		`"index": %d,` +
		`"timestamp": "%s",` +
		`"data": "%s",` +
		`"hash": %x` +
	`}`, block.index, block.timestamp.Format(time.RFC3339Nano), block.data, block.hash)
}
