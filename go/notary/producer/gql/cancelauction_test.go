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

func TestCancelAuctionStorageReturnsError(t *testing.T) {
	counter := 0

	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	inAuction := storage.Auction{
		ID:     "123abc",
		Seller: crypto.PubkeyToAddress(alice.PublicKey).Hex(),
	}

	mock := mockup.Tx{
		AuctionFunc:       func(ID string) (*storage.Auction, error) { return &inAuction, nil },
		AuctionCancelFunc: func(ID string) error { return errors.New("error") },
		RollbackFunc:      func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelAuctionInput{}
	in.AuctionId = "123abc"
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelAuction(struct{ Input input.CancelAuctionInput }{in})
	assert.Error(t, err, "error")
	assert.Equal(t, counter, 1)

}

func TestCancelAuctionStorageReturnsErrorOnWrongSigner(t *testing.T) {
	counter := 0

	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	bob, _ := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	inAuction := storage.Auction{
		ID:     "123abc",
		Seller: crypto.PubkeyToAddress(alice.PublicKey).Hex(),
	}

	mock := mockup.Tx{
		AuctionFunc:       func(ID string) (*storage.Auction, error) { return &inAuction, nil },
		AuctionCancelFunc: func(ID string) error { return errors.New("error") },
		RollbackFunc:      func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelAuctionInput{}
	in.AuctionId = "123abc"
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bob)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelAuction(struct{ Input input.CancelAuctionInput }{in})
	assert.Error(t, err, "Signer of CancelAuction is not the Seller")
	assert.Equal(t, counter, 0)
}

func TestCancelAuctionStorageReturnsOK(t *testing.T) {
	counter := 0

	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	inAuction := storage.Auction{
		ID:     "123abc",
		Seller: crypto.PubkeyToAddress(alice.PublicKey).Hex(),
	}

	mock := mockup.Tx{
		AuctionFunc:       func(ID string) (*storage.Auction, error) { return &inAuction, nil },
		AuctionCancelFunc: func(ID string) error { return nil },
		CommitFunc:        func() error { return nil },
		RollbackFunc:      func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelAuctionInput{}
	in.AuctionId = "123abc"
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelAuction(struct{ Input input.CancelAuctionInput }{in})
	assert.NilError(t, err)
}

// func TestCancelAuctionStorageReturnsErrNonExistingAuction(t *testing.T) {
// 	counter := 0

// 	alice, _ := crypto.HexToECDSA("4B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

// 	inAuction := storage.Auction{
// 		ID:     "123abc",
// 		Seller: crypto.PubkeyToAddress(alice.PublicKey).Hex(),
// 	}

// 	mock := mockup.Tx{
// 		AuctionFunc:       func(ID string) (*storage.Auction, error) { return &inAuction, nil },
// 		AuctionCancelFunc: func(ID string) error { return nil },
// 		RollbackFunc:      func() error { counter++; return nil },
// 	}
// 	service := &mockup.StorageService{
// 		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
// 	}

// 	in := input.CancelAuctionInput{}
// 	in.AuctionId = "123abc"
// 	hash, err := in.Hash()
// 	assert.NilError(t, err)
// 	signature, err := signer.Sign(hash.Bytes(), alice)
// 	assert.NilError(t, err)
// 	in.Signature = hex.EncodeToString(signature)

// 	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
// 	_, err = r.CancelAuction(struct{ Input input.CancelAuctionInput }{in})
// 	assert.NilError(t, err)
// }
