package entity

import (
	"crypto/sha256"
	"encoding/hex"
)

type RawTransaction struct {
	TransactionData string `json: "transactionData"`
	MaxGasFee       string `json: "maxGasFee"`
}

func (t *RawTransaction) Hash() (string, error) {
	transactionBin := make([]byte, hex.DecodedLen(len(t.TransactionData)))
	hex.Decode(transactionBin, []byte(t.TransactionData))
	hash := sha256.Sum256(transactionBin)
	hash = sha256.Sum256(hash[:])
	txId := hex.EncodeToString(reverse(hash[:]))
	return txId, nil
}

func reverse(input []byte) []byte {
	l := len(input)
	reversed := make([]byte, l)
	for i, n := range input {
		j := l - i - 1
		reversed[j] = n
	}
	return reversed
}
