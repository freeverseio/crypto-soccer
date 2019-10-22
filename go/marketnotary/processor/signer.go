package processor

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts/market"
)

type Signer struct {
	assets *market.Market
}

func NewSigner(marketContract *market.Market) *Signer {
	return &Signer{marketContract}
}

func (b *Signer) RSV(signature string) (r [32]byte, s [32]byte, v uint8, err error) {
	if len(signature) != 132 {
		return r, s, v, errors.New("wrong signature length")
	}
	signature = signature[2:] // remove 0x
	vect, err := hex.DecodeString(signature[0:64])
	if err != nil {
		return r, s, v, err
	}
	copy(r[:], vect)
	vect, err = hex.DecodeString(signature[64:128])
	if err != nil {
		return r, s, v, err
	}
	copy(s[:], vect)
	vect, err = hex.DecodeString(signature[128:130])
	v = vect[0]
	return r, s, v, err
}

func (b *Signer) HashPrivateMsg(currencyId uint8, price *big.Int, rnd *big.Int) ([32]byte, error) {
	privateHash, err := b.assets.HashPrivateMsg(
		&bind.CallOpts{},
		currencyId,
		price,
		rnd,
	)
	return privateHash, err
}

func (b *Signer) HashSellMessage(currencyId uint8, price *big.Int, rnd *big.Int, validUntil *big.Int, playerId *big.Int) ([32]byte, error) {
	var hash [32]byte
	hashPrivateMessage, err := b.assets.HashPrivateMsg(
		&bind.CallOpts{},
		currencyId,
		price,
		rnd,
	)
	if err != nil {
		return hash, err
	}
	hash, err = b.assets.BuildPutAssetForSaleTxMsg(
		&bind.CallOpts{},
		hashPrivateMessage,
		validUntil,
		playerId,
	)
	if err != nil {
		return hash, err
	}
	hash, err = b.assets.Prefixed(&bind.CallOpts{}, hash)
	return hash, err
}

// func (b *Signer) HashBuyMessage(currencyId uint8, price *big.Int, rnd *big.Int, validUntil *big.Int, playerId *big.Int, teamId *big.Int) ([32]byte, error) {
// 	var hash [32]byte
// 	hashPrivateMessage, err := b.assets.HashPrivateMsg(
// 		&bind.CallOpts{},
// 		currencyId,
// 		price,
// 		rnd,
// 	)
// 	if err != nil {
// 		return hash, err
// 	}
// 	sellMsgHash, err := b.assets.BuildPutForSaleTxMsg(
// 		&bind.CallOpts{},
// 		hashPrivateMessage,
// 		validUntil,
// 		playerId,
// 	)
// 	if err != nil {
// 		return hash, err
// 	}
// 	prefixedHash, err := b.assets.Prefixed(&bind.CallOpts{}, sellMsgHash)
// 	if err != nil {
// 		return hash, err
// 	}
// 	hash, err = b.assets.BuildAgreeToBuyTxMsg(
// 		&bind.CallOpts{},
// 		prefixedHash,
// 		teamId,
// 	)
// 	if err != nil {
// 		return hash, err
// 	}
// 	hash, err = b.assets.Prefixed(&bind.CallOpts{}, hash)
// 	return hash, err
// }
