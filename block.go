package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
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
