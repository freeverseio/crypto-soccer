package v1

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"fmt"
)

const sandboxURL = "https://api-sandbox.truust.io/1.0"
const sandboxPublicKey = "pk_stage_ZkNpNElWeEg="

type MarketPayContext struct {
}

func (c MarketPayContext) GetEndPoint() string {
	return "https://api-sandbox.truust.io/1.0"
}

func (c MarketPayContext) GetPublicKey() string {
	return "pk_stage_ZkNpNElWeEg="
}

type MarketPay struct {
	endpoint  string
	publicKey string
}

func NewMarketPay(context MarketPayContext) (*MarketPay, error) {
	return &MarketPay{
		context.GetEndPoint(),
		context.GetPublicKey(),
	}, nil
}

func New() (*MarketPay, error) {
	return NewMarketPay(MarketPayContext{})
}

func (b *MarketPay) CreateOrder(
	name string,
	value string,
) (*Order, error) {
	url := b.endpoint + "/express"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("name", name)
	_ = writer.WriteField("amount", value)
	_ = writer.WriteField("source", b.publicKey)
	_ = writer.WriteField("trustee_confirmed_url", "https://freeverseio.github.io/buyer/success")
	_ = writer.WriteField("trustee_denied_url", "https://freeverseio.github.io/buyer/failure")
	_ = writer.WriteField("settlor_confirmed_url", "https://freeverseio.github.io/seller/success")
	_ = writer.WriteField("settlor_denied_url", "https://freeverseio.github.io/seller/failure")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
	}

	// req.Header.Add("Accept", "application/json")
	// req.Header.Add("Authorization", b.bearerToken)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	// fmt.Println(string(body))
	order := &Order{}
	err = json.Unmarshal(body, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (b *MarketPay) GetOrder(hash string) (*Order, error) {
	url := b.endpoint + "/express/hash/" + hash
	method := "GET"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	if err != nil {
		fmt.Println(err)
	}
	// req.Header.Add("Accept", "application/json")
	// req.Header.Add("Authorization", b.bearerToken)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	// fmt.Println(string(body))
	order := &Order{}
	err = json.Unmarshal(body, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (b *MarketPay) IsPaid(order Order) bool {
	return order.Status == "PUBLISHED"
}
