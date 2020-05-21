package playstore

import (
	"context"

	"github.com/awa/go-iap/playstore"
	"google.golang.org/api/androidpublisher/v3"
)

type GoogleClientService struct {
	client *playstore.Client
}

func NewGoogleClientService(credentials []byte) (*GoogleClientService, error) {
	client, err := playstore.New(credentials)
	if err != nil {
		return nil, err
	}
	return &GoogleClientService{
		client: client,
	}, nil
}

func (b GoogleClientService) GetPurchase(
	ctx context.Context,
	packageName string,
	productId string,
	purchaseToken string,
) (*androidpublisher.ProductPurchase, error) {
	return b.client.VerifyProduct(
		ctx,
		packageName,
		productId,
		purchaseToken,
	)
}

func (b GoogleClientService) AcknowledgedPurchase(
	ctx context.Context,
	packageName string,
	productId string,
	purchaseToken string,
	payload string,
) error {
	return b.client.AcknowledgeProduct(
		ctx,
		packageName,
		productId,
		purchaseToken,
		payload,
	)
}
