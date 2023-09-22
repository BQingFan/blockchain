package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

// The dificulty of mining
const targetBits = 24

// Max nonce to avoid overflow
const maxNonce = math.MaxInt64

/*
Block structure with basic information.
*/
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	// nonce is required to verify a proof
	Nonce int
}

/*
Calculating the hashed value for a block.
Take block fieldsm concatenate them, and calculate a SHA-256 hash
on the concatenated combination.
*/
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

/*
Building a new block and return the pointer.
*/
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)

	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash[:]
	return block
}

/*
Blockchain structure as an ordered, back-linked list.
*/
type Blockchain struct {
	blocks []*Block
}

/*
Build the first block in chain.
*/
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

/*
Add blocks to the blockchain
*/
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

/*
Create a blockchain
*/
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

/*
Proof-of-work algorithum: Hashcash
First, take some publicly known data(block headers).
Then add a counter to it. The counter starts at 0.
Get a hash of the data+counter combination.
Check that the hash meets certain requirements.
If it does, we are done, else we need to increase the counter and repeat.
*/
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	// Left out the first targetBits 0s
	// Target is the biggest number meets the requirement.
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{b, target}
}

/*
Prepare the data to be checked.
*/
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{pow.block.PrevBlockHash, pow.block.Data, IntToBytes(pow.block.Timestamp), IntToBytes(int64(targetBits)), IntToBytes(int64(nonce))}, []byte{})
	return data
}

func IntToBytes(n int64) []byte {
	// Create a buffer to hold the binary representation
	buf := make([]byte, binary.MaxVarintLen64)

	// Encode the integer into the buffer
	binary.PutVarint(buf, n)

	// Trim any unused bytes from the buffer and return it
	return buf[:binary.Size(n)]
}

/*
Run the proof-of-work algorithm on the block to be checked.
*/
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		// prepare data
		data := pow.prepareData(nonce)
		// hash it with SHA-256
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		// convert the hash to a big integer
		hashInt.SetBytes((hash[:]))
		// compare the integer with the target
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Print("\n\n")

	// return iteration number nonce, and the hash it has
	return nonce, hash[:]
}

/*
Validate prood of works
Check for the stored nonce, valid if the hash is smaller than the stored target
*/
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}

/*
Serialize a block to store it in the database (BoltDB).
Keys and values are both byte arrays in BlotDB.
*/
func (b *Block) Serialize() []byte {
	// declare a buffer that will store the serialized data
	var result bytes.Buffer
	// initialize a gob encoder
	encoder := gob.NewEncoder(&result)
	// encode the block
	err := encoder.Encode(b)
	return result.Bytes()
}

/*
Deserialize a byte array to a Block.
*/
func DeserializeBlock(d []byte) *Block {
	// declare a block to store the deserialized data
	var block Block
	// initialize a gob decoder
	decoder := gob.NewDecoder(bytes.NewReader(d))
	// decode the block
	err := decoder.Decode(&block)
	return &block
}

func main() {
	bc := NewBlockchain()
	bc.AddBlock("Send 1 BTC to Amy")
	bc.AddBlock("Send 2 more BTC to Amy")

	for _, block := range bc.blocks {
		fmt.Printf("Prev.hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
	}
}
