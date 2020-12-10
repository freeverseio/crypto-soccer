package purchasevoider

import (
	"context"

	"google.golang.org/api/androidpublisher/v3"
)

type VoidPurchasesService interface {
	VoidedPurchases(ctx context.Context) ([]*androidpublisher.VoidedPurchase, error)
}
