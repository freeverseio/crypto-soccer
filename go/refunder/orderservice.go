package refunder

import (
	"context"

	"google.golang.org/api/androidpublisher/v3"
)

type OrderService interface {
	List(context.Context) *androidpublisher.PurchasesVoidedpurchasesListCall
}
