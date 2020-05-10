package staker_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/staker"
	"gotest.tools/assert"
)

func TestStakerNew(t *testing.T) {
	pvc, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	assert.NilError(t, err)
	s, err := staker.New(pvc)
	assert.NilError(t, err)
	assert.Equal(t, s.Address().Hex(), "0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7")
}

func TestStakerIsEnrolled(t *testing.T) {
	pvc, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	assert.NilError(t, err)
	s, err := staker.New(pvc)
	assert.NilError(t, err)
	isEnrolled, err := s.IsEnrolled(*bc.Contracts)
	assert.NilError(t, err)
	assert.Assert(t, !isEnrolled)
}

func TestStakerEnrolled(t *testing.T) {
	pvc, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	assert.NilError(t, err)
	s, err := staker.New(pvc)
	assert.NilError(t, err)
	assert.NilError(t, s.Enroll(*bc.Contracts))
}
