package helper

import (
	"crypto/ecdsa"
	"encoding/hex"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	Key *ecdsa.PrivateKey
}

func NewAccount() *Account {
	key, err := crypto.GenerateKey()
	if err != nil {
		log.Panic(err)
	}
	return &Account{key}
}

func (b Account) Address() common.Address {
	return crypto.PubkeyToAddress(b.Key.PublicKey)
}

func (b Account) PrivateKey() string {
	return hex.EncodeToString(b.Key.D.Bytes())
}
