package playstore

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	androidpublisher "google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

type Client struct {
}

type ClientService interface {
	GetPurchase(ctx context.Context, packageName string, productId string, purchaseToken string) (*androidpublisher.ProductPurchase, error)
	AcknowledgePurchase(ctx context.Context, packageName string, productId string, purchaseToken string, payload string) error
	Refund(ctx context.Context, packageName string, orderId string) error
}

type GoogleClientService struct {
	service *androidpublisher.Service
}

func NewGoogleClientService(credentials []byte) (*GoogleClientService, error) {
	c := &http.Client{Timeout: 10 * time.Second}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, c)

	conf, err := google.JWTConfigFromJSON(credentials, androidpublisher.AndroidpublisherScope)
	if err != nil {
		return nil, err
	}

	val := conf.Client(ctx).Transport.(*oauth2.Transport)
	_, err = val.Source.Token()
	if err != nil {
		return nil, err
	}

	service, err := androidpublisher.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		return nil, err
	}

	return &GoogleClientService{service}, nil
}

func (b GoogleClientService) Refund(
	ctx context.Context,
	packageName string,
	orderId string,
) error {
	ps := androidpublisher.NewOrdersService(b.service)
	return ps.Refund(packageName, orderId).Context(ctx).Do()
}

func (b GoogleClientService) GetPurchase(
	ctx context.Context,
	packageName string,
	productId string,
	purchaseToken string,
) (*androidpublisher.ProductPurchase, error) {
	ps := androidpublisher.NewPurchasesProductsService(b.service)
	return ps.Get(packageName, productId, purchaseToken).Context(ctx).Do()
}

func (b GoogleClientService) AcknowledgePurchase(
	ctx context.Context,
	packageName string,
	productId string,
	purchaseToken string,
	payload string,
) error {
	ps := androidpublisher.NewPurchasesProductsService(b.service)
	acknowledgeRequest := &androidpublisher.ProductPurchasesAcknowledgeRequest{DeveloperPayload: payload}
	err := ps.Acknowledge(packageName, productId, purchaseToken, acknowledgeRequest).Context(ctx).Do()
	return err
}
