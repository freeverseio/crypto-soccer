package gql_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/mockup"
	"gotest.tools/assert"
)

func TestGetAuctionPass(t *testing.T) {
	ch := make(chan interface{}, 10)

	mock := mockup.Tx{
		AuctionPassFunc: func(owner string) (*storage.AuctionPass, error) {
			ap := storage.NewAuctionPass()
			ap.Owner = "yo"
			return ap, nil
		},
		RollbackFunc: func() error { return nil },
		CommitFunc:   func() error { return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)

	in := input.HasAuctionPassInput{}
	in.Owner = "274877906944"

	auctionPass, err := r.HasAuctionPass(struct{ Input input.HasAuctionPassInput }{in})
	assert.NilError(t, err)
	assert.Equal(t, *auctionPass, false)

	in.Owner = "yo"

	auctionPass, err = r.HasAuctionPass(struct{ Input input.HasAuctionPassInput }{in})
	assert.NilError(t, err)
	assert.Equal(t, *auctionPass, true)
}

func TestGetNilAuctionPass(t *testing.T) {
	ch := make(chan interface{}, 10)

	mock := mockup.Tx{
		AuctionPassFunc: func(owner string) (*storage.AuctionPass, error) { return nil, nil },
		RollbackFunc:    func() error { return nil },
		CommitFunc:      func() error { return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)

	in := input.HasAuctionPassInput{}
	in.Owner = "274877906944"

	_, err := r.HasAuctionPass(struct{ Input input.HasAuctionPassInput }{in})
	assert.NilError(t, err)

}
