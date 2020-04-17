package input

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/graph-gophers/graphql-go"
)

type CancelAuctionInput struct {
	Signature string
	AuctionId graphql.ID
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
	if len(sign) != 65 {
		return false, fmt.Errorf("signature must be 65 bytes long")
	}
	if sign[64] != 27 && sign[64] != 28 {
		return false, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sign[64] -= 27 // Transform yellow paper V from 27/28 to 0/
	return signer.VerifySignature(hash[:], sign)
}

func (b CancelAuctionInput) SignerAddress() (common.Address, error) {
	hash := b.Hash()
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return common.Address{}, err
	}
	return signer.AddressFromSignature(hash.Bytes(), sign)
}
