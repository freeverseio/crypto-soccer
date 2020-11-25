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

	in := input.GetAuctionPassInput{}
	in.Signature = "a67621b4763db406f404c4a600ce0e79ee50147c209e85d2f146f0d760c0a1ac2a213a06f702995cee279af1f588b55c9fa462b2e6a9502d25cede77ec690ced1c"
	in.TeamId = "274877906944"

	auctionPass, err := r.GetAuctionPass(struct{ Input input.GetAuctionPassInput }{in})
	assert.NilError(t, err)
	assert.Assert(t, auctionPass.Owner() == "yo")
}
