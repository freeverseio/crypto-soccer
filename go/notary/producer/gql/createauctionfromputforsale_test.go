package gql_test

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/mockup"
	"gotest.tools/assert"
)

func TestCreateAuctionCallRollbackOnError(t *testing.T) {
	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)
	// We will here assign the next available team to offerer so she can make an offer for the players of a different team
	// we choose that player as a player very far from the current amount of teams (x2)
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, timezoneIdx, countryIdx)
	offererTeamIdx := nHumanTeams.Int64()
	sellerTeamIdx := offererTeamIdx + 1
	// offererTeamId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, big.NewInt(offererTeamIdx))
	offerer, _ := crypto.HexToECDSA("9B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	seller, _ := crypto.HexToECDSA("0A878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	playerId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, (big.NewInt(2 + 18*sellerTeamIdx)))

	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(offerer.PublicKey),
	)
	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(seller.PublicKey),
	)

	rollbackCounter := 0

	mock := mockup.Tx{
		AuctionInsertFunc: func(auction storage.Auction) error { return errors.New("errorAuctionInsertFunc") },
		RollbackFunc:      func() error { rollbackCounter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CreatePutPlayerForSaleInput{}
	in.ValidUntil = strconv.FormatInt(time.Now().Unix()+1000, 10)
	in.PlayerId = playerId.String()
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321

	hash, err := in.SellerDigest()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), seller)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CreateAuctionFromPutForSale(struct {
		Input input.CreatePutPlayerForSaleInput
	}{in})
	assert.Error(t, err, "errorAuctionInsertFunc")
	assert.Equal(t, rollbackCounter, 1)
}

func TestCreateAuctionCallCommit(t *testing.T) {
	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)
	// We will here assign the next available team to offerer so she can make an offer for the players of a different team
	// we choose that player as a player very far from the current amount of teams (x2)
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, timezoneIdx, countryIdx)
	offererTeamIdx := nHumanTeams.Int64()
	sellerTeamIdx := offererTeamIdx + 1
	// offererTeamId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, big.NewInt(offererTeamIdx))
	offerer, _ := crypto.HexToECDSA("9B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	seller, _ := crypto.HexToECDSA("0A878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	playerId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, (big.NewInt(2 + 18*sellerTeamIdx)))

	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(offerer.PublicKey),
	)
	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(seller.PublicKey),
	)

	counter := 0

	mock := mockup.Tx{
		AuctionInsertFunc: func(auction storage.Auction) error { return nil },
		CommitFunc:        func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CreatePutPlayerForSaleInput{}
	in.ValidUntil = strconv.FormatInt(time.Now().Unix()+100, 10)
	in.PlayerId = playerId.String()
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	hash, err := in.SellerDigest()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), seller)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	result, err := r.CreateAuctionFromPutForSale(struct {
		Input input.CreatePutPlayerForSaleInput
	}{in})
	assert.NilError(t, err)
	assert.Equal(t, counter, 1)

	// test that the return coincides with the ID of the auction
	id, err := in.ID()
	assert.NilError(t, err)
	assert.Equal(t, id, result)
}
