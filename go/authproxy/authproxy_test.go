package authproxy

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

const (
	pvkhex = "c6d398e89bf7cbda7663ca881bd992eb80ad170e4ca0bd65a8b1c719ee02bc67"
)

func TestBasic(t *testing.T) {
	pvk, _ := crypto.HexToECDSA(pvkhex)
	addr := crypto.PubkeyToAddress(pvk.PublicKey)

	token, err := SignToken(pvk, time.Now())
	assert.Nil(t, err)
	addr1, _, err := VerifyToken(token, time.Hour)
	assert.Nil(t, err)

	assert.EqualValues(t, addr, addr1)
}

func TestHalfGraceUp(t *testing.T) {
	pvk, _ := crypto.HexToECDSA(pvkhex)

	token, err := SignToken(pvk, time.Now().Add(time.Hour/2))
	assert.Nil(t, err)
	_, _, err = VerifyToken(token, time.Hour)
	assert.Nil(t, err)
}

func TestHalfGraceDown(t *testing.T) {
	pvk, _ := crypto.HexToECDSA(pvkhex)

	token, err := SignToken(pvk, time.Now().Add(-time.Hour/2))
	assert.Nil(t, err)
	_, _, err = VerifyToken(token, time.Hour)
	assert.Nil(t, err)
}

func TestDobuleGraceUp(t *testing.T) {
	pvk, _ := crypto.HexToECDSA(pvkhex)

	token, err := SignToken(pvk, time.Now().Add(time.Hour*2))
	assert.Nil(t, err)
	_, _, err = VerifyToken(token, time.Hour)
	assert.NotNil(t, err)
}

func TestDobuleGraceDown(t *testing.T) {
	pvk, _ := crypto.HexToECDSA(pvkhex)

	token, err := SignToken(pvk, time.Now().Add(-time.Hour*2))
	assert.Nil(t, err)
	_, _, err = VerifyToken(token, time.Hour)
	assert.NotNil(t, err)
}

func BenchmarkSignatureVerification(b *testing.B) {
	pvkhex := "c6d398e89bf7cbda7663ca881bd992eb80ad170e4ca0bd65a8b1c719ee02bc67"
	pvk, _ := crypto.HexToECDSA(pvkhex)

	token, err := SignToken(pvk, time.Now())
	assert.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyToken(token, time.Hour)
	}
}
