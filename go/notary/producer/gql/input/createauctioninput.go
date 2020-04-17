package input

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/graph-gophers/graphql-go"
)

type CreateAuctionInput struct {
	Signature  string
	PlayerId   string
	CurrencyId int32
	Price      int32
	Rnd        int32
	ValidUntil string
}

func (b CreateAuctionInput) ID() (graphql.ID, error) {
	hash, err := b.Hash()
	if err != nil {
		return graphql.ID(""), err
	}
	return graphql.ID(hash.String()[2:]), nil
}

func (b CreateAuctionInput) Hash() (common.Hash, error) {
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return common.Hash{}, errors.New("invalid playerId")
	}
	validUntil, err := strconv.ParseInt(b.ValidUntil, 10, 64)
	if err != nil {
		return common.Hash{}, err
	}
	hash, err := signer.HashSellMessage(
		uint8(b.CurrencyId),
		big.NewInt(int64(b.Price)),
		big.NewInt(int64(b.Rnd)),
		validUntil,
		playerId,
	)
	return hash, err
}

func (b CreateAuctionInput) VerifySignature() (bool, error) {
	hash, err := b.Hash()
	if err != nil {
		return false, err
	}
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
	sign[64] -= 27 // Transform yellow paper V from 27/28 to 0/1
	return signer.VerifySignature(hash.Bytes(), sign)
}

func (b CreateAuctionInput) SignerAddress() (common.Address, error) {
	hash, err := b.Hash()
	if err != nil {
		return common.Address{}, err
	}
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return common.Address{}, err
	}
	if len(sign) != 65 {
		return common.Address{}, fmt.Errorf("signature must be 65 bytes long")
	}
	if sign[64] != 27 && sign[64] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sign[64] -= 27 // Transform yellow paper V from 27/28 to 0/1
	return signer.AddressFromSignature(hash.Bytes(), sign)
}
