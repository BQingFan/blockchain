package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

/*
Block structure with basic information.
*/
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
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
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()
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

func main() {
	bc := NewBlockchain()
	bc.AddBlock("Send 1 BTC to Amy")
	bc.AddBlock("Send 2 more BTC to Amy")

	for _, block := range bc.blocks {
		fmt.Printf("Prev.hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}
}
