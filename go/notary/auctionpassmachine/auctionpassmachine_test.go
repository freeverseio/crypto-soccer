package auctionpassmachine_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/auctionpassmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestAuctionPassMachineCreation(t *testing.T) {
	client := NewMockClientService()
	iapTestOn := true
	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()
	order := storage.NewAuctionPassPlaystoreOrder()
	_, err = auctionpassmachine.New(
		tx,
		client,
		*order,
		*bc.Contracts,
		bc.Owner,
		iapTestOn,
	)
	assert.NilError(t, err)
}

func TestMachineCreationFailedState(t *testing.T) {
	client := NewMockClientService()
	iapTestOn := true
	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()
	order := storage.NewAuctionPassPlaystoreOrder()
	order.State = storage.AuctionPassPlaystoreOrderFailed
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
}

func TestMachineCreationRefundedState(t *testing.T) {
	client := NewMockClientService()
	iapTestOn := true
	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()
	order := storage.NewAuctionPassPlaystoreOrder()
	order.State = storage.AuctionPassPlaystoreOrderRefunded
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
}

func TestMachineCreationCompleteState(t *testing.T) {
	client := NewMockClientService()
	iapTestOn := true
	tx, err := service.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()
	order := storage.NewAuctionPassPlaystoreOrder()
	order.State = storage.AuctionPassPlaystoreOrderComplete
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
}
