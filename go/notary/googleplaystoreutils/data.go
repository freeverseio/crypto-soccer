package googleplaystoreutils

import (
	"encoding/json"
)

type Data struct {
	OrderId       string
	PackageName   string
	ProductId     string
	PurchaseToken string
}

func DataFromReceipt(receipt string) (*Data, error) {
	var temp0 struct{ Payload string }
	if err := json.Unmarshal([]byte(receipt), &temp0); err != nil {
		return nil, err
	}

	var temp1 struct{ Json string }
	if err := json.Unmarshal([]byte(temp0.Payload), &temp1); err != nil {
		return nil, err
	}

	data := Data{}
	if err := json.Unmarshal([]byte(temp1.Json), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (b Data) ToReceipt() string {
	return ""
}
