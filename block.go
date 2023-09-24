package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

/*
Block structure with basic information.
*/
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

/*
Calculating the hashed value for the transactions in the block.
*/
func (b *Block) HashTansactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

/*
Building a new block and return the pointer.
*/
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)

	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash[:]
	return block
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
	if err != nil {
		log.Panic(err)
	}
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
	if err != nil {
		log.Panic(err)
	}
	return &block
}

/*
Build the first block in chain.
*/
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}
