package input_test

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
	"gotest.tools/assert"
)

func TestSetTeamNameInputHash(t *testing.T) {
	in := input.SetTeamNameInput{}
	hash, err := in.Hash()
	assert.Error(t, err, "Invalid TeamId")
	// s := hex.EncodeToString("0x01")
	// fmt.Println("peruba 0x01", s)
	in.TeamId = "3"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x074b4277787bca36334cf57f0507141ef743a08d7690dba02af123626e6955d0")

	in.Name = "ciao"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x70203f52f1e52e6727239d5197821e1ed77161d7a04a45ad1fb15702c69433d4")
}

func TestSetTeamNameInputSign(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	assert.NilError(t, err)

	in := input.SetTeamNameInput{}
	in.TeamId = "4"
	in.Name = "ciao"

	hash, err := in.Hash()
	assert.NilError(t, err)
	hash = helper.PrefixedHash(hash)
	//assert.Equal(t, "sdfasdfa", hash.Hex())
	sign, err := helper.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "3feac668bb718f492638b9b58d1f294379cdc8bde40074f5e49c3f80f28190e121f0fd08227c64a643dd032748ef772b0d1cf1500f649345521c133290c941a91b")

	isValid, err := helper.VerifySignature(hash, sign)
	assert.NilError(t, err)
	assert.Assert(t, isValid)

	sender, err := in.SignerAddress()
	assert.NilError(t, err)
	assert.Equal(t, sender.Hex(), "0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7")
}
