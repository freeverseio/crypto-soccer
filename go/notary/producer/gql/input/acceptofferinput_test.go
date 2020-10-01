package input_test

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"gotest.tools/assert"
)

func TestAcceptOfferInputHash(t *testing.T) {
	in := input.AcceptOfferInput{}
	in.ValidUntil = "2000000000"
	in.OfferValidUntil = "199999000"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	hash, err := in.SellerDigest()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0xe30150bc666d8a20396c30794d1d1eaf86ecf427d4f4e3d8a4aa87d4aa3fc4b5")
}
func TestAcceptOfferInputAuctionID(t *testing.T) {
	in := input.AcceptOfferInput{}
	in.ValidUntil = "2000000000"
	in.OfferValidUntil = "199999000"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	id, err := in.AuctionID()
	assert.NilError(t, err)
	assert.Equal(t, string(id), "4fa64a72d14d52385a359bb69326e4faca26771ccb0662cb6fa5555fdf97beb1")
}
func TestAcceptOfferSignerAddress(t *testing.T) {
	in := input.AcceptOfferInput{}
	in.ValidUntil = "2000000000"
	in.OfferValidUntil = "199999000"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321

	hash, err := in.SellerDigest()
	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	signature, err := signer.Sign(hash.Bytes(), alice)

	in.Signature = hex.EncodeToString(signature)
	address, err := in.SignerAddress()
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), crypto.PubkeyToAddress(alice.PublicKey).Hex())
}

func TestAcceptOfferIsSignerOwner(t *testing.T) {
	tz := uint8(1)
	countryIdxInTz := big.NewInt(0)
	// We will here assign the next available team to alice so she can put a playerId for sale
	// playerId from the second team is made an offer
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, tz, countryIdxInTz)
	nPlayerInCountryForAlice := nHumanTeams.Int64()*18 + 1
	playerId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, tz, countryIdxInTz, big.NewInt(nPlayerInCountryForAlice))
	alice, _ := crypto.HexToECDSA("5B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	in := input.AcceptOfferInput{}
	in.ValidUntil = strconv.FormatInt(time.Now().Unix()+1000, 10)
	in.OfferValidUntil = strconv.FormatInt(time.Now().Unix()+100, 10)
	in.PlayerId = playerId.String()
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321

	hash, err := in.SellerDigest()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	isOwner, err := in.IsSignerOwnerOfPlayer(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, isOwner, false)

	isValid, err := in.IsValidForBlockchainFreeze(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, isValid, false)

	tx, err := bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		tz,
		countryIdxInTz,
		crypto.PubkeyToAddress(alice.PublicKey),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = helper.WaitReceipt(bc.Client, tx, 5)
	if err != nil {
		t.Fatal(err)
	}

	isOwner, err = in.IsSignerOwnerOfPlayer(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, isOwner, true)

	isValid, err = in.IsValidForBlockchainFreeze(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, isValid, true)
}
