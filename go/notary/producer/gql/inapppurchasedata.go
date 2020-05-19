package gql

type InappPurchaseData struct {
	OrderId       string
	PackageName   string
	ProductId     string
	PurchaseToken string
}

func InappPurchaseDataFromReceipt(receipt string) (*InappPurchaseData, error) {
	data := InappPurchaseData{}
	return &data, nil
}
