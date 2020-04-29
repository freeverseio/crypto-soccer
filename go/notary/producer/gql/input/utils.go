package input

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
)

func VerifySignature(hash common.Hash, sign []byte) (bool, error) {
	if len(sign) != 65 {
		return false, fmt.Errorf("signature must be 65 bytes long")
	}
	if sign[64] != 27 && sign[64] != 28 {
		return false, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sg := make([]byte, len(sign))
	copy(sg, sign)
	sg[64] -= 27 // Transform yellow paper V from 27/28 to 0/1
	return signer.VerifySignature(hash.Bytes(), sg)
}

func AddressFromSignature(hash common.Hash, sign []byte) (common.Address, error) {
	if len(sign) != 65 {
		return common.Address{}, fmt.Errorf("signature must be 65 bytes long")
	}
	if sign[64] != 27 && sign[64] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sg := make([]byte, len(sign))
	copy(sg, sign)
	sg[64] -= 27 // Transform yellow paper V from 27/28 to 0/
	return signer.AddressFromSignature(hash.Bytes(), sg)
}

func Sign(hash []byte, pvr *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := crypto.Sign(hash, pvr)
	if len(sig) != 65 {
		return []byte{}, fmt.Errorf("signature must be 65 bytes long")
	}
	if sig[64] != 0 && sig[64] != 1 {
		return []byte{}, fmt.Errorf("invalid Ethereum signature (V is not 0 or 1)")
	}
	sig[64] += 27
	return sig, err
}
