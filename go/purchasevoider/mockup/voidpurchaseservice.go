package mockup

import (
	"context"

	"google.golang.org/api/androidpublisher/v3"
)

type VoidPurchaseService struct {
	VoidedPurchasesFn func(context.Context) ([]*androidpublisher.VoidedPurchase, error)
}

func (b VoidPurchaseService) VoidedPurchases(ctx context.Context) ([]*androidpublisher.VoidedPurchase, error) {
	return b.VoidedPurchasesFn(ctx)
}
