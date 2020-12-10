package mockup

import (
	"context"

	"google.golang.org/api/androidpublisher/v3"
)

type VoidPurchasesService struct {
	VoidedPurchasesFn func(context.Context) ([]*androidpublisher.VoidedPurchase, error)
}

func (b VoidPurchasesService) VoidedPurchases(ctx context.Context) ([]*androidpublisher.VoidedPurchase, error) {
	return b.VoidedPurchasesFn(ctx)
}
