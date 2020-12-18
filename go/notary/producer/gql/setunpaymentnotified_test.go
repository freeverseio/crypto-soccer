package gql_test

import (
	"errors"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/mockup"
	"gotest.tools/assert"
)

func TestSetUnpaymentNotified(t *testing.T) {
	ch := make(chan interface{}, 10)

	mock := mockup.Tx{
		UnpaymentUpdateNotifiedFunc: func(unpayment storage.Unpayment) error { return nil },
		RollbackFunc:                func() error { return nil },
		CommitFunc:                  func() error { return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)

	in := input.SetUnpaymentNotifiedInput{}
	in.Id = "274877906944"

	_, err := r.SetUnpaymentNotified(struct {
		Input input.SetUnpaymentNotifiedInput
	}{in})
	assert.NilError(t, err)
}

func TestSetUnpaymentNotifiedErr(t *testing.T) {
	ch := make(chan interface{}, 10)

	mock := mockup.Tx{
		UnpaymentUpdateNotifiedFunc: func(unpayment storage.Unpayment) error { return errors.New("Error") },
		RollbackFunc:                func() error { return nil },
		CommitFunc:                  func() error { return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)

	in := input.SetUnpaymentNotifiedInput{}
	in.Id = "274877906944"

	_, err := r.SetUnpaymentNotified(struct {
		Input input.SetUnpaymentNotifiedInput
	}{in})
	assert.Error(t, err, "Error")
}
