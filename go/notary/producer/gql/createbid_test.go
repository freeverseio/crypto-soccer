package gql_test

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/mockup"

	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"gotest.tools/assert"
)

func TestCreateBidStorageReturnError(t *testing.T) {
	counter := 0

	mock := mockup.Tx{
		BidInsertFunc: func(bid storage.Bid) error { return errors.New("error") },
		RollbackFunc:  func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)

	in := input.CreateBidInput{}
	in.AuctionId = "unexistent"
	in.ExtraPrice = 332
	in.Rnd = 1243523
	in.TeamId = "274877906945"
	in.Signature = "4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1c"

	err := r.CreateBid(struct{ Input input.CreateBidInput }{in})
	assert.Error(t, err, "signer is not the owner of teamId 274877906945")
}

func TestCreateBidStorageReturnsOK(t *testing.T) {
	counter := 0

	mock := mockup.Tx{
		BidInsertFunc: func(bid storage.Bid) error { return nil },
		CommitFunc:    func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)

	in := input.CreateBidInput{}
	in.AuctionId = "yeah"
	in.ExtraPrice = 332
	in.Rnd = 1243523
	in.TeamId = "274877906945"
	in.Signature = "4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1c"

	err := r.CreateBid(struct{ Input input.CreateBidInput }{in})
	assert.Error(t, err, "signer is not the owner of teamId 274877906945")
}
