package input_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/graph-gophers/graphql-go"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestCreateBidInputHash(t *testing.T) {
	putforsale := input.CreatePutPlayerForSaleInput{}
	putforsale.ValidUntil = "2000000000"
	putforsale.PlayerId = "274877906944"
	putforsale.CurrencyId = 1
	putforsale.Price = 41234
	putforsale.Rnd = 42321
	auctionId, err := putforsale.ID()
	assert.NilError(t, err)
	assert.Equal(t, string(auctionId), "58912aae76687d592fefbb46a6192474eb56ce15eb12dea2e41bee3b9fca45d3")

	in := input.CreateBidInput{}
	in.AuctionId = auctionId
	in.ExtraPrice = 332
	in.Rnd = 1243523
	in.TeamId = "274877906945"

	hash, err := in.Hash(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0xaa3b6bde688ad7a6cb75ec987cf8b295dc6fd0782defe0da43f2c6a088e99c95")
}

func TestCreateBidInputSign(t *testing.T) {
	auction := storage.NewAuction()
	auction.ValidUntil = int64(2000000000)
	auction.OfferValidUntil = int64(1999999000)
	auction.PlayerID = "274877906944"
	auction.CurrencyID = 1
	auction.Price = 41234
	auction.Rnd = 42321

	auctionId, err := auction.ComputeID()
	assert.Equal(t, auctionId, "24f45caee25883ab36cc32e2b152b94b8a05bdb086ac68784e3a4686f4d961e8")

	in := input.CreateBidInput{}
	in.AuctionId = graphql.ID(auctionId)
	in.ExtraPrice = 332
	in.Rnd = 1243523
	in.TeamId = "274877906945"
	in.Signature = "4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b"

	isValid, err := in.IsSignerOwnerOfTeam(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, isValid, false)
}

func TestCreateGoodBidInputSign(t *testing.T) {
	tz := uint8(1)
	countryIdxInTz := big.NewInt(0)
	// We will here assign the next available team to alice so she can put a playerId for sale
	// playerId from the second team is made an offer
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, tz, countryIdxInTz)
	teamId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, tz, countryIdxInTz, nHumanTeams)
	alice, _ := crypto.HexToECDSA("0B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	in := input.CreateBidInput{}
	in.AuctionId = graphql.ID("32123")
	in.ExtraPrice = 332
	in.Rnd = 1243523
	in.TeamId = teamId.String()

	hash, err := in.Hash(*bc.Contracts)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	isValid, err := in.IsSignerOwnerOfTeam(*bc.Contracts)
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

	isValid, err = in.IsSignerOwnerOfTeam(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, isValid, true)
}
