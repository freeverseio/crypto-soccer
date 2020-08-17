package helper

import (
	"crypto/ecdsa"
	"encoding/hex"
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
	if err != nil {
		return nil, err
	}
	if len(sig) != 65 {
		return nil, fmt.Errorf("signature must be 65 bytes long")
	}
	if sig[64] != 0 && sig[64] != 1 {
		return nil, fmt.Errorf("invalid Ethereum signature (V is not 0 or 1)")
	}
	sig[64] += 27
	return sig, err
}

func PrefixedHash(hash common.Hash) common.Hash {
	ss := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(hash), hash)
	return crypto.Keccak256Hash([]byte(ss))
}

func RSV(signature string) (r [32]byte, s [32]byte, v uint8, err error) {
	if len(signature) != 132 && len(signature) != 130 {
		return r, s, v, fmt.Errorf("wrong signature length %v", len(signature))
	}
	if len(signature) == 132 {
		signature = signature[2:] // remove 0x
	}
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
