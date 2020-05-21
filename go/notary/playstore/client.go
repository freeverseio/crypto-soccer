package playstore

import (
	"context"

	"google.golang.org/api/androidpublisher/v3"
)

type Client struct {
}

type ClientService interface {
	GetPurchase(ctx context.Context, packageName string, productId string, purchaseToken string) (*androidpublisher.ProductPurchase, error)
	AcknowledgedPurchase(ctx context.Context, packageName string, productId string, purchaseToken string, payload string) error
}
