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

type VoidPurchasesService struct {
	service     *androidpublisher.Service
	packageName string
}

func NewVoidPurchasesService(
	credentials []byte,
	packageName string,
) (refunder.VoidPurchasesService, error) {
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

	return &VoidPurchasesService{
		service,
		packageName,
	}, nil
}

func (b VoidPurchasesService) VoidedPurchases(ctx context.Context) ([]*androidpublisher.VoidedPurchase, error) {
	ps := androidpublisher.NewPurchasesVoidedpurchasesService(b.service)
	list := ps.List(b.packageName)
	response, err := list.Context(ctx).Do()
	return response.VoidedPurchases, err
}
