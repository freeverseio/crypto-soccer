package gql

import log "github.com/sirupsen/logrus"

type InappPurchaseData struct {
	OrderId       string
	PackageName   string
	ProductId     string
	PurchaseToken string
}

func InappPurchaseDataFromReceipt(receipt string) (*InappPurchaseData, error) {
	log.Info(receipt)
	data := InappPurchaseData{}
	return &data, nil
}
