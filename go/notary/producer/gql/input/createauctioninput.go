package input

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
)

type CreateAuctionInput struct {
	Signature  string
	PlayerId   string
	CurrencyId int
	Price      int
	Rnd        int
	ValidUntil string
}

func (b CreateAuctionInput) Hash() (common.Hash, error) {
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	validUntil, err := strconv.ParseInt(b.ValidUntil, 10, 64)
	if err != nil {
		return [32]byte{}, err
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
	return signer.VerifySignature(hash[:], []byte(b.Signature))
}

func (b CreateAuctionInput) SignerAddress() (common.Address, error) {
	hash, err := b.Hash()
	if err != nil {
		return common.Address{}, err
	}
	return signer.AddressFromSignature(hash[:], []byte(b.Signature))
}
