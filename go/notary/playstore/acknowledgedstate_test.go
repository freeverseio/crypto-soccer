package playstore_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/googleplaystoreutils"
	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"google.golang.org/api/androidpublisher/v3"
	"gotest.tools/assert"
)

func TestAckStateProcessErrorInCLient(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	order.State = storage.PlaystoreOrderAcknowledged
	client := googleplaystoreutils.NewMockClientService()
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
	assert.Equal(t, m.Order().State, storage.PlaystoreOrderAcknowledged)
	assert.Equal(t, m.Order().StateExtra, "not implemented")
}

func TestAckStateProcessTestPurchaseWithTestOff(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	order.State = storage.PlaystoreOrderAcknowledged
	client := googleplaystoreutils.NewMockClientService()
	client.GetPurchaseFunc = func() (*androidpublisher.ProductPurchase, error) {
		return &androidpublisher.ProductPurchase{
			PurchaseType: new(int64),
		}, nil
	}
	iapTestOn := false
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
	assert.Equal(t, m.Order().StateExtra, "test order")
}

func TestAckStateProcessTestPurchaseWithTestOn(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	order.State = storage.PlaystoreOrderAcknowledged
	client := googleplaystoreutils.NewMockClientService()
	client.GetPurchaseFunc = func() (*androidpublisher.ProductPurchase, error) {
		return &androidpublisher.ProductPurchase{
			PurchaseType: new(int64),
		}, nil
	}
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
	assert.Equal(t, m.Order().StateExtra, "invalid player")
}

func TestAckStateProcess(t *testing.T) {
	order := storage.NewPlaystoreOrder()
	order.State = storage.PlaystoreOrderAcknowledged
	client := googleplaystoreutils.NewMockClientService()
	client.GetPurchaseFunc = func() (*androidpublisher.ProductPurchase, error) {
		return &androidpublisher.ProductPurchase{}, nil
	}
	iapTestOn := false
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
	assert.Equal(t, m.Order().StateExtra, "invalid player")
}
