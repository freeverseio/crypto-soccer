package auctionpassmachine_test

import (
	"errors"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/auctionpassmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestRefundingStateProcessError(t *testing.T) {
	order := storage.NewAuctionPassPlaystoreOrder()
	order.State = storage.AuctionPassPlaystoreOrderRefunding
	client := NewMockClientService()
	client.RefundFunc = func() error { return errors.New("complicated") }
	iapTestOn := true
	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()
	m, err := auctionpassmachine.New(
		tx,
		client,
		*order,
		*bc.Contracts,
		bc.Owner,
		iapTestOn,
	)
	assert.NilError(t, err)
	assert.NilError(t, m.Process())
	assert.Equal(t, m.Order().State, storage.AuctionPassPlaystoreOrderRefunding)
	assert.Equal(t, m.Order().StateExtra, "complicated")
}

func TestRefundingStateProcessNoError(t *testing.T) {
	order := storage.NewAuctionPassPlaystoreOrder()
	order.State = storage.AuctionPassPlaystoreOrderRefunding
	client := NewMockClientService()
	client.RefundFunc = func() error { return nil }
	iapTestOn := true
	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()
	m, err := auctionpassmachine.New(
		tx,
		client,
		*order,
		*bc.Contracts,
		bc.Owner,
		iapTestOn,
	)
	assert.NilError(t, err)
	assert.NilError(t, m.Process())
	assert.Equal(t, m.Order().State, storage.AuctionPassPlaystoreOrderRefunded)
	assert.Equal(t, m.Order().StateExtra, "")
}
