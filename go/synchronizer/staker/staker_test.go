package staker_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/helper"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/staker"
	"gotest.tools/assert"
)

func TestStakerNew(t *testing.T) {
	account := helper.NewAccount()

	s, err := staker.New(account.Key)
	assert.NilError(t, err)
	assert.Equal(t, s.Address().Hex(), account.Address().Hex())
}

func TestStakerIsTrustedParty(t *testing.T) {
	account := helper.NewAccount()

	s, err := staker.New(account.Key)
	assert.NilError(t, err)
	isTrusted, err := s.IsTrustedParty(*bc.Contracts)
	assert.NilError(t, err)
	assert.Assert(t, !isTrusted)
}

func TestStakerIsEnrolled(t *testing.T) {
	account := helper.NewAccount()

	s, err := staker.New(account.Key)

	isEnrolled, err := s.IsEnrolled(*bc.Contracts)
	assert.NilError(t, err)
	assert.Assert(t, !isEnrolled)
}

func TestStakerSubmitRoot(t *testing.T) {
	account := helper.NewAccount()

	s, err := staker.New(account.Key)
	assert.NilError(t, err)

	assert.Error(t, s.Init(*bc.Contracts), "[staker] not a trusted party")
	root := [32]byte{0x0}
	assert.Error(t, s.SubmitRoot(*bc.Contracts, root), "failed to estimate gas needed: The execution failed due to an exception.")
}

func TestStakerEnroll(t *testing.T) {
	pvc, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	assert.NilError(t, err)

	s, err := staker.New(pvc)
	assert.NilError(t, err)

	assert.Error(t, s.Init(*bc.Contracts), "[staker] not a trusted party")

	t.Run("be trusted party", func(t *testing.T) {
		tx, err := bc.Contracts.Stakers.AddTrustedParty(bind.NewKeyedTransactor(bc.Owner), s.Address())
		assert.NilError(t, err)
		_, err = helper.WaitReceipt(bc.Client, tx, 60)
		assert.NilError(t, err)
		isTrusted, err := s.IsTrustedParty(*bc.Contracts)
		assert.NilError(t, err)
		assert.Assert(t, isTrusted)
	})

	assert.NilError(t, s.Init(*bc.Contracts))
	root := [32]byte{0x0}
	assert.NilError(t, s.SubmitRoot(*bc.Contracts, root)) // TODO should fail
}
