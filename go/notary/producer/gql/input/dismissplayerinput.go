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
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

type DismissPlayerInput struct {
	Signature       string
	PlayerId        string
	ValidUntil      string
	ReturnToAcademy bool
}

func (b DismissPlayerInput) Hash() (common.Hash, error) {
	uint256Ty, _ := abi.NewType("int256", "int256", nil)

	arguments := abi.Arguments{
		{Type: uint256Ty},
		{Type: uint256Ty},
	}

	validUntil, _ := new(big.Int).SetString(b.ValidUntil, 10)
	if validUntil == nil {
		return common.Hash{}, errors.New("invalid validUntil")
	}
	playerId, _ := new(big.Int).SetString(string(b.PlayerId), 10)
	if playerId == nil {
		return common.Hash{}, errors.New("invalid playerId")
	}

	bytes, err := arguments.Pack(
		validUntil,
		playerId,
	)
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(bytes), nil
}

func (b DismissPlayerInput) SignerAddress() (common.Address, error) {
	hash, err := b.Hash()
	if err != nil {
		return common.Address{}, err
	}
	hash = helper.PrefixedHash(hash)
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return common.Address{}, err
	}
	return helper.AddressFromHashAndSignature(hash, sign)
}

func (b DismissPlayerInput) IsSignerOwner(contracts contracts.Contracts) (bool, error) {
	signerAddress, err := b.SignerAddress()
	if err != nil {
		return false, err
	}
	playerId, _ := new(big.Int).SetString(string(b.PlayerId), 10)
	if playerId == nil {
		return false, errors.New("invalid teamId")
	}
	owner, err := contracts.Market.GetOwnerPlayer(&bind.CallOpts{}, playerId)
	fmt.Printf("Owner %v", owner.Hex())
	fmt.Printf("Signer address %v", signerAddress.Hex())
	if err != nil {
		return false, err
	}
	return signerAddress == owner, nil
}

func (b DismissPlayerInput) IsPlayerOnSale(tx storage.Tx) (bool, string, error) {
	var auctionID = ""
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return false, auctionID, errors.New("invalid playerId")
	}

	auctions, err := tx.AuctionsByPlayerId(b.PlayerId)
	if err != nil {
		return false, auctionID, err
	}

	isOnSale := false
	for _, auction := range auctions {
		if (auction.State != storage.AuctionCancelled) && (auction.State != storage.AuctionEnded) && (auction.State != storage.AuctionFailed) {
			isOnSale = true
			auctionID = auction.ID
		}
	}
	return isOnSale, auctionID, nil
}
