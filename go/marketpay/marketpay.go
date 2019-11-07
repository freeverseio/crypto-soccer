package marketpay

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type Customer struct {
	ID int `json:"id"`
}

type MarketPay struct {
	endpoint string
}

func New() (*MarketPay, error) {
	return &MarketPay{
		"https://api-sandbox.truust.io",
	}, nil
}

func (b *MarketPay) CreateCustomer(
	prefix string,
	phoneNumber string,
) (float64, error) {
	url := b.endpoint + "/2.0/customers"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("prefix", prefix)
	_ = writer.WriteField("phone", phoneNumber)
	err := writer.Close()
	if err != nil {
		return 0, err
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return 0, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer sk_stage_P4aQCVzTRGhub2p4k2Fl6YxQ")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	// fmt.Println(string(body))
	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}
	customerID := dat["data"].(map[string]interface{})["id"].(float64)
	return customerID, nil
}
