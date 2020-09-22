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

func TestCancelOfferStorageReturnsError1(t *testing.T) {
	counter := 0

	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	inOffer := storage.Offer{
		AuctionID: "123abc",
		Seller:    crypto.PubkeyToAddress(alice.PublicKey).Hex(),
		State:     storage.OfferStarted,
	}

	mock := mockup.Tx{
		OfferFunc:       func(AuctionID string) (*storage.Offer, error) { return &inOffer, nil },
		OfferCancelFunc: func(AuctionID string) error { return errors.New("error") },
		RollbackFunc:    func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelOfferInput{}
	in.OfferId = "123abc"
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelOffer(struct{ Input input.CancelOfferInput }{in})
	assert.Error(t, err, "error")
	assert.Equal(t, counter, 1)

}

func TestCancelOfferStorageReturnsErrorOnWrongSigner(t *testing.T) {
	counter := 0

	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	bob, _ := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	inOffer := storage.Offer{
		AuctionID: "123abc",
		Seller:    crypto.PubkeyToAddress(alice.PublicKey).Hex(),
		State:     storage.OfferStarted,
	}

	mock := mockup.Tx{
		OfferFunc:       func(AuctionID string) (*storage.Offer, error) { return &inOffer, nil },
		OfferCancelFunc: func(AuctionID string) error { return errors.New("error") },
		RollbackFunc:    func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelOfferInput{}
	in.OfferId = "123abc"
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bob)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelOffer(struct{ Input input.CancelOfferInput }{in})
	assert.Error(t, err, "Signer of Canceloffer is neither the Seller nor the Buyer")
	assert.Equal(t, counter, 1)
}

func TestCancelOfferStorageReturnsOK(t *testing.T) {
	counter := 0

	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	inOffer := storage.Offer{
		AuctionID: "123abc",
		Seller:    crypto.PubkeyToAddress(alice.PublicKey).Hex(),
		State:     storage.OfferStarted,
	}

	mock := mockup.Tx{
		OfferFunc:       func(AuctionID string) (*storage.Offer, error) { return &inOffer, nil },
		OfferCancelFunc: func(AuctionID string) error { return nil },
		CommitFunc:      func() error { return nil },
		RollbackFunc:    func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelOfferInput{}
	in.OfferId = "123abc"
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelOffer(struct{ Input input.CancelOfferInput }{in})
	assert.NilError(t, err)
}

func TestCancelOfferStorageReturnsErrOnEmptyOfferId(t *testing.T) {
	counter := 0

	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	inOffer := storage.Offer{
		AuctionID: "123abc",
		Seller:    crypto.PubkeyToAddress(alice.PublicKey).Hex(),
		State:     storage.OfferStarted,
	}

	mock := mockup.Tx{
		OfferFunc:       func(AuctionID string) (*storage.Offer, error) { return &inOffer, nil },
		OfferCancelFunc: func(AuctionID string) error { return nil },
		CommitFunc:      func() error { return nil },
		RollbackFunc:    func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelOfferInput{}
	in.OfferId = ""
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelOffer(struct{ Input input.CancelOfferInput }{in})
	assert.Error(t, err, "empty OfferId when trying to cancel an offer")
}

func TestCancelOfferStorageReturnsErrOnWrongOfferId(t *testing.T) {
	counter := 0

	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	mock := mockup.Tx{
		OfferFunc:       func(AuctionID string) (*storage.Offer, error) { return nil, nil },
		OfferCancelFunc: func(AuctionID string) error { return nil },
		CommitFunc:      func() error { return nil },
		RollbackFunc:    func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelOfferInput{}
	in.OfferId = "3245786876a1"
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelOffer(struct{ Input input.CancelOfferInput }{in})
	assert.Error(t, err, "could not find an offer to cancel")
}

func TestCancelOfferStorageReturnsErrorOnWronfOfferState(t *testing.T) {
	counter := 0

	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	inOffer := storage.Offer{
		AuctionID: "123abc",
		Seller:    crypto.PubkeyToAddress(alice.PublicKey).Hex(),
		State:     storage.OfferAccepted,
	}

	mock := mockup.Tx{
		OfferFunc:       func(AuctionID string) (*storage.Offer, error) { return &inOffer, nil },
		OfferCancelFunc: func(AuctionID string) error { return nil },
		CommitFunc:      func() error { return nil },
		RollbackFunc:    func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelOfferInput{}
	in.OfferId = "123abc"
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelOffer(struct{ Input input.CancelOfferInput }{in})
	assert.Error(t, err, "cannot cancel an offer unless it is in Started state")
}
