package input

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/graph-gophers/graphql-go"
)

type CancelAuctionInput struct {
	Signature string
	AuctionId graphql.ID
}

func (b CancelAuctionInput) Hash() (common.Hash, error) {
	bytes32Ty, _ := abi.NewType("bytes32", "bytes32", nil)
	arguments := abi.Arguments{
		{
			Type: bytes32Ty,
		},
	}
	idArray := [32]byte{}
	auctionHex, err := hex.DecodeString(string(b.AuctionId))
	if err != nil {
		return common.Hash{}, err
	}
	copy(idArray[:], auctionHex)
	bytes, err := arguments.Pack(idArray)
	if err != nil {
		return common.Hash{}, err
	}
	hash := [32]byte{}
	copy(hash[:], crypto.Keccak256Hash(bytes).Bytes())
	ss := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(hash), hash)
	return crypto.Keccak256Hash([]byte(ss)), nil
}

func (b CancelAuctionInput) SignerAddress() (common.Address, error) {
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
