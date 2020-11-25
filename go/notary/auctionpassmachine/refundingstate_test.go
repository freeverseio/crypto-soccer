package auctionpassmachine_test

import (
	"errors"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestRefundingStateProcessError(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	order.State = storage.PlaystoreOrderRefunding
	client := NewMockClientService()
	client.RefundFunc = func() error { return errors.New("complicated") }
	iapTestOn := true
	m, err := playstore.New(
		client,
		*order,
		*bc.Contracts,
		bc.Owner,
		namesdb,
		iapTestOn,
	)
	assert.NilError(t, err)
	assert.NilError(t, m.Process())
	assert.Equal(t, m.Order().State, storage.PlaystoreOrderRefunding)
	assert.Equal(t, m.Order().StateExtra, "complicated")
}

func TestRefundingStateProcessNoError(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	order.State = storage.PlaystoreOrderRefunding
	client := NewMockClientService()
	client.RefundFunc = func() error { return nil }
	iapTestOn := true
	m, err := playstore.New(
		client,
		*order,
		*bc.Contracts,
		bc.Owner,
		namesdb,
		iapTestOn,
	)
	assert.NilError(t, err)
	assert.NilError(t, m.Process())
	assert.Equal(t, m.Order().State, storage.PlaystoreOrderRefunded)
	assert.Equal(t, m.Order().StateExtra, "")
}
