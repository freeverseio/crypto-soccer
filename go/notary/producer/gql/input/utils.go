package input

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
)

func verifySignature(hash common.Hash, sign []byte) (bool, error) {
	if len(sign) != 65 {
		return false, fmt.Errorf("signature must be 65 bytes long")
	}
	if sign[64] != 27 && sign[64] != 28 {
		return false, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sign[64] -= 27 // Transform yellow paper V from 27/28 to 0/1
	return signer.VerifySignature(hash.Bytes(), sign)
}

func addressFromSignature(hash common.Hash, sign []byte) (common.Address, error) {
	if len(sign) != 65 {
		return common.Address{}, fmt.Errorf("signature must be 65 bytes long")
	}
	if sign[64] != 27 && sign[64] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sign[64] -= 27 // Transform yellow paper V from 27/28 to 0/1
	return signer.AddressFromSignature(hash.Bytes(), sign)
}
