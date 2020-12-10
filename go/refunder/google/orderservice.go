package google

import (
	"context"
	"net/http"
	"time"

	"github.com/freeverseio/crypto-soccer/go/refunder"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

type OrderService struct {
	service     *androidpublisher.Service
	packageName string
}

func NewOrderService(
	credentials []byte,
	packageName string,
) (refunder.OrderService, error) {
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

	return &OrderService{
		service,
		packageName,
	}, nil
}

func (b OrderService) List(ctx context.Context) *androidpublisher.PurchasesVoidedpurchasesListCall {
	ps := androidpublisher.NewPurchasesVoidedpurchasesService(b.service)
	list := ps.List(b.packageName)
	return list
}
