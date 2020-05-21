package playstore

import "google.golang.org/api/androidpublisher/v3"

type PurchaseValidator struct {
	purchase androidpublisher.ProductPurchase
}

func NewPurchaseValidator(purchase androidpublisher.ProductPurchase) *PurchaseValidator {
	return &PurchaseValidator{
		purchase: purchase,
	}
}

func (b PurchaseValidator) IsPurchased() bool {
	return b.purchase.PurchaseState == 0
}

func (b PurchaseValidator) IsCanceled() bool {
	return b.purchase.PurchaseState == 1
}

func (b PurchaseValidator) IsPending() bool {
	return b.purchase.PurchaseState == 2
}

func (b PurchaseValidator) IsAcknowledged() bool {
	return b.purchase.AcknowledgementState == 1
}
