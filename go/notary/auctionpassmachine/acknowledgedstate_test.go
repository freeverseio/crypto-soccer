package auctionpassmachine_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/auctionpassmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/googleplaystoreutils"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/androidpublisher/v3"
	"gotest.tools/assert"
)

func TestAckStateProcessErrorInCLient(t *testing.T) {
	order := storage.NewAuctionPassPlaystoreOrder()
	order.State = storage.AuctionPassPlaystoreOrderAcknowledged
	client := googleplaystoreutils.NewMockClientService()
	iapTestOn := true
	log.Infof("Hola")
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
	assert.Equal(t, m.Order().State, storage.AuctionPassPlaystoreOrderAcknowledged)
	assert.Equal(t, m.Order().StateExtra, "not implemented")
}

func TestAckStateProcessTestPurchaseWithTestOff(t *testing.T) {
	order := storage.NewAuctionPassPlaystoreOrder()
	order.State = storage.AuctionPassPlaystoreOrderAcknowledged
	client := googleplaystoreutils.NewMockClientService()
	client.GetPurchaseFunc = func() (*androidpublisher.ProductPurchase, error) {
		return &androidpublisher.ProductPurchase{
			PurchaseType: new(int64),
		}, nil
	}
	iapTestOn := false
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
	assert.Equal(t, m.Order().State, storage.AuctionPassPlaystoreOrderComplete)
	assert.Equal(t, m.Order().StateExtra, "test order")
}

func TestAckStateProcessTestPurchaseWithTestOn(t *testing.T) {
	order := storage.NewAuctionPassPlaystoreOrder()
	order.State = storage.AuctionPassPlaystoreOrderAcknowledged
	client := googleplaystoreutils.NewMockClientService()
	client.GetPurchaseFunc = func() (*androidpublisher.ProductPurchase, error) {
		return &androidpublisher.ProductPurchase{
			PurchaseType: new(int64),
		}, nil
	}
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
	assert.Equal(t, m.Order().StateExtra, "invalid team")
}

func TestAckStateProcess(t *testing.T) {
	order := storage.NewAuctionPassPlaystoreOrder()
	order.State = storage.AuctionPassPlaystoreOrderAcknowledged
	client := googleplaystoreutils.NewMockClientService()
	client.GetPurchaseFunc = func() (*androidpublisher.ProductPurchase, error) {
		return &androidpublisher.ProductPurchase{}, nil
	}
	iapTestOn := false
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
	assert.Equal(t, m.Order().StateExtra, "invalid team")
}
