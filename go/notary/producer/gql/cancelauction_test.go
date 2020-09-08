package gql_test

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/mockup"
	"gotest.tools/assert"
)

func TestCancelAuctionStorageReturnsError(t *testing.T) {
	counter := 0

	mock := mockup.Tx{
		AuctionCancelFunc: func(ID string) error { return errors.New("error") },
		RollbackFunc:      func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelAuctionInput{}
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelAuction(struct{ Input input.CancelAuctionInput }{in})
	assert.Error(t, err, "error")
	assert.Equal(t, counter, 1)
}

func TestCancelAuctionStorageReturnsOK(t *testing.T) {
	counter := 0

	mock := mockup.Tx{
		AuctionCancelFunc: func(ID string) error { return nil },
		CommitFunc:        func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CancelAuctionInput{}
	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CancelAuction(struct{ Input input.CancelAuctionInput }{in})
	assert.NilError(t, err)
	assert.Equal(t, counter, 1)
}
