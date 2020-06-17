package playstore_test

import (
	"errors"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"google.golang.org/api/androidpublisher/v3"
	"gotest.tools/assert"
)

func TestOpenStateProcessInvalidPurchaseState(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	client := NewMockClientService()
	client.GetPurchaseFunc = func() (*androidpublisher.ProductPurchase, error) {
		return &androidpublisher.ProductPurchase{
			PurchaseState: 3,
		}, nil
	}
	client.AcknowledgedPurchaseFunc = func() error { return nil }
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
	assert.Equal(t, m.Order().State, storage.PlaystoreOrderFailed)
	assert.Equal(t, m.Order().StateExtra, "invalid puchase state")
}

func TestOpenStateProcessPendingPurchaseState(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	client := NewMockClientService()
	client.GetPurchaseFunc = func() (*androidpublisher.ProductPurchase, error) {
		return &androidpublisher.ProductPurchase{
			PurchaseState: 2,
		}, nil
	}
	client.AcknowledgedPurchaseFunc = func() error { return nil }
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
	assert.Equal(t, m.Order().State, storage.PlaystoreOrderOpen)
	assert.Equal(t, m.Order().StateExtra, "pending")
}

func TestOpenStateProcessCancelledPurchaseState(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	client := NewMockClientService()
	client.GetPurchaseFunc = func() (*androidpublisher.ProductPurchase, error) {
		return &androidpublisher.ProductPurchase{
			PurchaseState: 1,
		}, nil
	}
	client.AcknowledgedPurchaseFunc = func() error { return nil }
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
	assert.Equal(t, m.Order().State, storage.PlaystoreOrderComplete)
	assert.Equal(t, m.Order().StateExtra, "cancelled")
}

func TestOpenStateProcessPurchasedPurchaseState(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	order.OrderId = "3"
	order.TeamId = "1099511627776"
	order.PlayerId = "333"
	client := NewMockClientService()
	client.GetPurchaseFunc = func() (*androidpublisher.ProductPurchase, error) {
		return &androidpublisher.ProductPurchase{
			PurchaseState: 0,
		}, nil
	}
	client.AcknowledgedPurchaseFunc = func() error { return nil }
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
	assert.Equal(t, m.Order().StateExtra, "orderId 3 has an invalid playerId 333")
}

func TestOpenStateProcessErrorInAck(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	client := NewMockClientService()
	client.GetPurchaseFunc = func() (*androidpublisher.ProductPurchase, error) {
		return &androidpublisher.ProductPurchase{
			PurchaseState: 0,
		}, nil
	}
	client.AcknowledgedPurchaseFunc = func() error { return errors.New("horrorrrrr") }
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
	assert.Equal(t, m.Order().State, storage.PlaystoreOrderOpen)
	assert.Equal(t, m.Order().StateExtra, "invalid teamId")
}

func TestOpenStateProcessAlreadyAck(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	client := NewMockClientService()
	client.GetPurchaseFunc = func() (*androidpublisher.ProductPurchase, error) {
		return &androidpublisher.ProductPurchase{
			PurchaseState:        0,
			AcknowledgementState: 1,
		}, nil
	}
	client.AcknowledgedPurchaseFunc = func() error { return errors.New("horrorrrrr") }
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
	assert.Equal(t, m.Order().State, storage.PlaystoreOrderFailed)
	assert.Equal(t, m.Order().StateExtra, "already acknowledged")
}
