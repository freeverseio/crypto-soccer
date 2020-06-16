package playstore_test

import (
	"context"
	"errors"

	"google.golang.org/api/androidpublisher/v3"
)

type MockClientService struct {
	GetPurchaseFunc          func() (*androidpublisher.ProductPurchase, error)
	AcknowledgedPurchaseFunc func() error
}

func NewMockClientService() *MockClientService {
	return &MockClientService{
		GetPurchaseFunc: func() (*androidpublisher.ProductPurchase, error) {
			return nil, errors.New("not implemented")
		},
		AcknowledgedPurchaseFunc: func() error {
			return errors.New("not implemented")
		},
	}
}

func (b MockClientService) GetPurchase(
	ctx context.Context,
	packageName string,
	productId string,
	purchaseToken string,
) (*androidpublisher.ProductPurchase, error) {
	return b.GetPurchaseFunc()
}

func (b MockClientService) AcknowledgePurchase(
	ctx context.Context,
	packageName string,
	productId string,
	purchaseToken string,
	payload string,
) error {
	return b.AcknowledgedPurchaseFunc()
}
