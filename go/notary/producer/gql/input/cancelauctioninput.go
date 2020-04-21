package input

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
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
	copy(idArray[:], []byte(b.AuctionId))
	bytes, err := arguments.Pack(idArray)
	if err != nil {
		return common.Hash{}, err
	}
	hash := [32]byte{}
	copy(hash[:], crypto.Keccak256Hash(bytes).Bytes())
	ss := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(hash), hash)
	return crypto.Keccak256Hash([]byte(ss)), nil
}

func (b CancelAuctionInput) VerifySignature() (bool, error) {
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
	sign[64] -= 27 // Transform yellow paper V from 27/28 to 0/
	return signer.VerifySignature(hash[:], sign)
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
	if len(sign) != 65 {
		return common.Address{}, fmt.Errorf("signature must be 65 bytes long")
	}
	if sign[64] != 27 && sign[64] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sign[64] -= 27 // Transform yellow paper V from 27/28 to 0/1
	return signer.AddressFromSignature(hash.Bytes(), sign)
}
