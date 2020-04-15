package input

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
)

type CreateAuctionInput struct {
	Signature  string
	PlayerId   string
	CurrencyId int32
	Price      int32
	Rnd        int32
	ValidUntil string
}

func (b CreateAuctionInput) ID() string {
	h := sha256.New()

	h.Write([]byte(fmt.Sprintf("%s%s%d%d%d%s", b.Signature, b.PlayerId, b.CurrencyId, b.Price, b.Rnd, b.ValidUntil)))
	return hex.EncodeToString(h.Sum(nil))
}

func (b CreateAuctionInput) Hash() (common.Hash, error) {
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
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
	return signer.VerifySignature(hash[:], sign)
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
	return signer.AddressFromSignature(hash[:], sign)
}
