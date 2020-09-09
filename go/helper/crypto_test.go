package helper_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/helper"
	"gotest.tools/assert"
)

func TestComputeSellPlayerDigest(t *testing.T) {
	currencyId := uint8(1)
	playerId := big.NewInt(11114324213423)
	price := big.NewInt(4324213423)
	rnd := big.NewInt(434324324213423)
	validUntil := uint32(235985749)
	auctionDurationAfterOfferIsAccepted := uint32(4358487)

	digest, err := helper.ComputeSellPlayerDigest(
		currencyId,
		price,
		rnd,
		validUntil,
		auctionDurationAfterOfferIsAccepted,
		playerId,
	)
	assert.NilError(t, err)
	assert.Equal(t, digest.Hex(), "0xf778fa056bd74980669505bf4666bbde172de50abe33d569f3ce597bdd81198b")

}
