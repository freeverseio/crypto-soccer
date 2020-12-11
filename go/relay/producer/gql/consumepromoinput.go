package gql

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
)

func (b ConsumePromoInput) Hash() (common.Hash, error) {
	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)

	arguments := abi.Arguments{
		{Type: uint256Ty},
		{Type: uint256Ty},
	}

	teamId, _ := new(big.Int).SetString(string(b.TeamId), 10)
	if teamId == nil {
		return common.Hash{}, errors.New("Invalid TeamId")
	}
	playerId, _ := new(big.Int).SetString(string(b.PlayerId), 10)
	if playerId == nil {
		return common.Hash{}, errors.New("Invalid PlayerId")
	}

	bytes, err := arguments.Pack(
		playerId,
		teamId,
	)
	if err != nil {
		return common.Hash{}, err
	}

	hash := [32]byte{}
	copy(hash[:], crypto.Keccak256Hash(bytes).Bytes())
	ss := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(hash), hash)
	return crypto.Keccak256Hash([]byte(ss)), nil
}

func (b ConsumePromoInput) SignerAddress() (common.Address, error) {
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

func (b ConsumePromoInput) IsSignerOwner(contracts contracts.Contracts) (bool, error) {
	signerAddress, err := b.SignerAddress()
	if err != nil {
		return false, err
	}
	teamId, _ := new(big.Int).SetString(string(b.TeamId), 10)
	if teamId == nil {
		return false, errors.New("Invalid teamId")
	}
	owner, err := contracts.Market.GetOwnerTeam(&bind.CallOpts{}, teamId)
	if err != nil {
		return false, err
	}
	return signerAddress == owner, nil
}
