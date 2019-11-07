package marketpay

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

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
) error {
	url := "https://api-sandbox.truust.io/2.0/customers"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("prefix", prefix)
	_ = writer.WriteField("phone", phoneNumber)
	err := writer.Close()
	if err != nil {
		return err
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer sk_stage_P4aQCVzTRGhub2p4k2Fl6YxQ")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}
