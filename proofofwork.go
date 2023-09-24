package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

// The dificulty of mining
const targetBits = 24

// Max nonce to avoid overflow
const maxNonce = math.MaxInt64

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
