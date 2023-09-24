package main

import (
	"github.com/boltdb/bolt"
)

// DB File to store blockchain information
const dbFile = "blockchain.db"

// Bucket name for storing blocks' information
const blocksBucket = "blocks"

/*
Blockchain structure.
Store the tip of the blockchain, and the db connection to the BoltDB.
*/
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

/*
Build the first block in chain.
*/
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

/*
Add blocks to the blockchain.
The new Block is stored on DB, and will work as the new tip of the blockchain.
*/
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	// find the last block in the blockchain
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	NewBlock := NewBlock(data, lastHash)

	// add the new block to db, update the tip in blockchain
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(NewBlock.Hash, NewBlock.Serialize())
		err = b.Put([]byte("l"), NewBlock.Hash)
		bc.tip = NewBlock.Hash

		return nil
	})
}

/*
Create a blockchain.
First open the DB file.
Secondly we need to check if there is a blockchain in the DB file.

If there is a blockchain: create a new blockchain instance,
set the tip of the blockchain to the last block hash stored in the DB.

If there is no blockchain: create the genesis block, store it in the DB,
save the genesis block's hash as the last block hash,
create a new blockchain instance with its tip pointing at the genesis block.
*/
func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}

// iterator
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := BlockchainIterator{bc.tip, bc.db}
	return &bci
}

func (bci *BlockchainIterator) Next() *Block {
	var block *Block

	err := bci.db.View(func(tx *bolt.Tx) error {
		bc := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bc.Get(bci.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})

	bci.currentHash = block.PrevBlockHash
	return block
}

func main() {
	bc := NewBlockchain()
	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()

}
