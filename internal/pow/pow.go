package pow

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

var (
	// ErrNonceNotFound is used when a nonce is not found
	ErrNonceNotFound = fmt.Errorf("nonce not found")
)

// Compute computes a wantNonce for a message with a given difficulty
func Compute(msg []byte, difficulty int) (int, []byte, error) {
	nonce := 0
	target := newTarget(difficulty)
	nonceBytes := make([]byte, 4)
	for ; nonce < math.MaxInt64; nonce++ {
		binary.BigEndian.PutUint32(nonceBytes, uint32(nonce))
		pop := sha256.Sum256(append(msg, nonceBytes...))
		res := new(big.Int).SetBytes(pop[:])
		if res.Cmp(target) < 0 {
			return nonce, pop[:], nil
		}
	}
	return 0, nil, ErrNonceNotFound
}

// Verify verifies a message with a given difficulty and nonce
func Verify(msg, hash []byte, difficulty, nonce int) bool {
	target := newTarget(difficulty)
	nonceBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(nonceBytes, uint32(nonce))
	msgHash := sha256.Sum256(append(msg, nonceBytes...))
	hashInt := new(big.Int).SetBytes(msgHash[:])
	return hashInt.Cmp(target) < 0 && bytes.Equal(msgHash[:], hash)
}

func newTarget(difficulty int) *big.Int {
	target := big.NewInt(1)
	return target.Lsh(target, uint(256-difficulty))
}
