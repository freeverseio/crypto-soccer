package gql_test

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/mockup"
	"github.com/graph-gophers/graphql-go"
	"gotest.tools/assert"
)

func TestCancelAllOffersBySeller(t *testing.T) {
	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, timezoneIdx, countryIdx)
	playerOwnerTeamIdx := nHumanTeams.Int64()
	playerOwner, _ := crypto.HexToECDSA("9B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	playerId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, (big.NewInt(2 + 18*playerOwnerTeamIdx)))

	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(playerOwner.PublicKey),
	)
	counter := 0

	mock := mockup.Tx{
		CancelAllOffersByPlayerIdFunc: func(playerId string) error { return nil },
		OfferCancelFunc:               func(AuctionID string) error { return errors.New("error") },
		RollbackFunc:                  func() error { counter++; return nil },
		CommitFunc:                    func() error { return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelAllOffersBySellerInput{}
	in.PlayerId = graphql.ID(playerId.String())
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), playerOwner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)
	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelAllOffersBySeller(struct {
		Input input.CancelAllOffersBySellerInput
	}{in})
	assert.NilError(t, err)
	assert.Equal(t, counter, 0)

}

func TestCancelAllOffersBySellerStorageReturnsErrorOnWrongSigner(t *testing.T) {
	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, timezoneIdx, countryIdx)
	playerOwnerTeamIdx := nHumanTeams.Int64()
	playerOwner, _ := crypto.HexToECDSA("9B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	playerId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, (big.NewInt(2 + 18*playerOwnerTeamIdx)))

	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(playerOwner.PublicKey),
	)
	counter := 0
	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	mock := mockup.Tx{
		CancelAllOffersByPlayerIdFunc: func(playerId string) error { return nil },
		OfferCancelFunc:               func(AuctionID string) error { return errors.New("error") },
		RollbackFunc:                  func() error { counter++; return nil },
		CommitFunc:                    func() error { return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelAllOffersBySellerInput{}
	in.PlayerId = graphql.ID(playerId.String())
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice) //Now someone else signs
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)
	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelAllOffersBySeller(struct {
		Input input.CancelAllOffersBySellerInput
	}{in})
	errMsg := "signer is not the owner of playerId " + playerId.String()
	assert.Error(t, err, errMsg)
	assert.Equal(t, counter, 0)

}

func TestCancelAllOffersBySellerStorageReturnsErrOnWrongSql(t *testing.T) {
	counter := 0

	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, timezoneIdx, countryIdx)
	playerOwnerTeamIdx := nHumanTeams.Int64()
	playerOwner, _ := crypto.HexToECDSA("9B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	playerId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, (big.NewInt(2 + 18*playerOwnerTeamIdx)))

	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(playerOwner.PublicKey),
	)

	mock := mockup.Tx{
		CancelAllOffersByPlayerIdFunc: func(playerId string) error { return errors.New("error") },
		OfferCancelFunc:               func(AuctionID string) error { return nil },
		CommitFunc:                    func() error { return nil },
		RollbackFunc:                  func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelAllOffersBySellerInput{}
	in.PlayerId = graphql.ID(playerId.String())
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), playerOwner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelAllOffersBySeller(struct {
		Input input.CancelAllOffersBySellerInput
	}{in})
	assert.Error(t, err, "error")
}
