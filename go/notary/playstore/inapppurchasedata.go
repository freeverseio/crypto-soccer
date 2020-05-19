package playstore

import (
	"encoding/json"
)

type InappPurchaseData struct {
	OrderId       string
	PackageName   string
	ProductId     string
	PurchaseToken string
}

func InappPurchaseDataFromReceipt(receipt string) (*InappPurchaseData, error) {
	var temp0 struct{ Payload string }
	if err := json.Unmarshal([]byte(receipt), &temp0); err != nil {
		return nil, err
	}

	var temp1 struct{ Json string }
	if err := json.Unmarshal([]byte(temp0.Payload), &temp1); err != nil {
		return nil, err
	}

	data := InappPurchaseData{}
	if err := json.Unmarshal([]byte(temp1.Json), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (b InappPurchaseData) ToReceipt() string {
	return ""
}
