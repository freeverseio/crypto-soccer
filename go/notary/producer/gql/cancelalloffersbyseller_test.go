package gql_test

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/mockup"
	"gotest.tools/assert"
)

func TestCancelAllOffersBySellerStorageReturnsError1(t *testing.T) {
	counter := 0

	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	mockOfferAlice1 := storage.Offer{
		AuctionID: "123abc",
		Seller:    crypto.PubkeyToAddress(alice.PublicKey).Hex(),
		State:     storage.OfferStarted,
	}

	mockOfferAlice2 := storage.Offer{
		AuctionID: "123abcd",
		Seller:    crypto.PubkeyToAddress(alice.PublicKey).Hex(),
		State:     storage.OfferStarted,
	}

	mockOffersByPlayerId := []storage.Offer{mockOfferAlice1, mockOfferAlice2}

	mock := mockup.Tx{
		OffersStartedByPlayerIdFunc: func(playerId string) ([]storage.Offer, error) { return mockOffersByPlayerId, nil },
		OfferCancelFunc:             func(AuctionID string) error { return errors.New("error") },
		RollbackFunc:                func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelAllOffersBySellerInput{}
	in.PlayerId = "274877906944"
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelAllOffersBySeller(struct {
		Input input.CancelAllOffersBySellerInput
	}{in})
	assert.Error(t, err, "error")
	assert.Equal(t, counter, 1)

}

func TestCancelAllOffersBySellerStorageReturnsErrorOnWrongSigner(t *testing.T) {
	counter := 0

	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	bob, _ := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	mockOfferAlice1 := storage.Offer{
		AuctionID: "123abc",
		Seller:    crypto.PubkeyToAddress(alice.PublicKey).Hex(),
		State:     storage.OfferStarted,
	}

	mockOfferAlice2 := storage.Offer{
		AuctionID: "123abcd",
		Seller:    crypto.PubkeyToAddress(bob.PublicKey).Hex(),
		State:     storage.OfferStarted,
	}

	mockOffersByPlayerId := []storage.Offer{mockOfferAlice1, mockOfferAlice2}

	mock := mockup.Tx{
		OffersStartedByPlayerIdFunc: func(playerId string) ([]storage.Offer, error) { return mockOffersByPlayerId, nil },
		OfferCancelFunc:             func(AuctionID string) error { return errors.New("error") },
		RollbackFunc:                func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelAllOffersBySellerInput{}
	in.PlayerId = "274877906944"
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bob)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelAllOffersBySeller(struct {
		Input input.CancelAllOffersBySellerInput
	}{in})
	assert.Error(t, err, "Signer of Cancelalloffersbyseller is not the Seller")
	assert.Equal(t, counter, 1)
}

func TestCancelAllOffersBySellerStorageReturnsOK(t *testing.T) {
	counter := 0

	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	mockOfferAlice1 := storage.Offer{
		AuctionID: "123abc",
		Seller:    crypto.PubkeyToAddress(alice.PublicKey).Hex(),
		State:     storage.OfferStarted,
	}

	mockOfferAlice2 := storage.Offer{
		AuctionID: "123abcd",
		Seller:    crypto.PubkeyToAddress(alice.PublicKey).Hex(),
		State:     storage.OfferStarted,
	}

	mockOffersByPlayerId := []storage.Offer{mockOfferAlice1, mockOfferAlice2}

	mock := mockup.Tx{
		OffersStartedByPlayerIdFunc: func(playerId string) ([]storage.Offer, error) { return mockOffersByPlayerId, nil },
		OfferCancelFunc:             func(AuctionID string) error { return nil },
		CommitFunc:                  func() error { return nil },
		RollbackFunc:                func() error { counter++; return nil },
	}

	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelAllOffersBySellerInput{}
	in.PlayerId = "274877906944"
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelAllOffersBySeller(struct {
		Input input.CancelAllOffersBySellerInput
	}{in})
	assert.NilError(t, err)
}

func TestCancelAllOffersBySellerStorageReturnsErrOnWrongOfferId(t *testing.T) {
	counter := 0

	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	mock := mockup.Tx{
		OffersStartedByPlayerIdFunc: func(playerId string) ([]storage.Offer, error) { return nil, nil },
		OfferCancelFunc:             func(AuctionID string) error { return nil },
		CommitFunc:                  func() error { return nil },
		RollbackFunc:                func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelAllOffersBySellerInput{}
	in.PlayerId = "274877906944"
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelAllOffersBySeller(struct {
		Input input.CancelAllOffersBySellerInput
	}{in})
	assert.Error(t, err, "could not find an existing offers to cancel")
}
