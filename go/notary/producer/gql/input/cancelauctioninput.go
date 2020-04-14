package input

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/graph-gophers/graphql-go"
)

type CancelAuctionInput struct {
	Signature string
	ID        graphql.ID
}

func (b CancelAuctionInput) Hash() common.Hash {
	return crypto.Keccak256Hash([]byte(b.Signature))
}

func (b CancelAuctionInput) VerifySignature() (bool, error) {
	hash := b.Hash()
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return false, err
	}
	return signer.VerifySignature(hash[:], sign)
}
