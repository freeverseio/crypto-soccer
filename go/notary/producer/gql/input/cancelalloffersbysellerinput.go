package input

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/graph-gophers/graphql-go"
)

type CancelAllOffersBySellerInput struct {
	Signature string
	PlayerId  graphql.ID
}

func (b CancelAllOffersBySellerInput) Hash() (common.Hash, error) {
	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)
	arguments := abi.Arguments{
		{Type: uint256Ty},
	}
	playerId, _ := new(big.Int).SetString(string(b.PlayerId), 10)
	if playerId == nil {
		return common.Hash{}, errors.New("invalid playerId")
	}
	bytes, err := arguments.Pack(playerId)
	if err != nil {
		return common.Hash{}, err
	}
	hash := [32]byte{}
	copy(hash[:], crypto.Keccak256Hash(bytes).Bytes())
	ss := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(hash), hash)
	return crypto.Keccak256Hash([]byte(ss)), nil
}

func (b CancelAllOffersBySellerInput) SignerAddress() (common.Address, error) {
	hash, err := b.Hash()
	if err != nil {
		return common.Address{}, err
	}
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return common.Address{}, err
	}
	return helper.AddressFromHashAndSignature(hash, sign)
}

func (b CancelAllOffersBySellerInput) IsSignerOwnerOfPlayer(contracts contracts.Contracts) (bool, error) {
	signerAddress, err := b.SignerAddress()
	if err != nil {
		return false, err
	}
	playerId, _ := new(big.Int).SetString(string(b.PlayerId), 10)
	if playerId == nil {
		return false, errors.New("invalid playerId")
	}
	owner, err := contracts.Market.GetOwnerPlayer(&bind.CallOpts{}, playerId)
	if err != nil {
		return false, err
	}
	return signerAddress == owner, nil
}
