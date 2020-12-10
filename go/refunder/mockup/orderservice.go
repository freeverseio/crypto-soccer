package mockup

import (
	"context"

	"google.golang.org/api/androidpublisher/v3"
)

type OrderService struct {
	ListFn func(context.Context) *androidpublisher.PurchasesVoidedpurchasesListCall
}

func (b OrderService) List(ctx context.Context) *androidpublisher.PurchasesVoidedpurchasesListCall {
	return b.ListFn(ctx)
}
