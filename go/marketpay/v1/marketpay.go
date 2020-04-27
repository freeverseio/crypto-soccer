package v1

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"fmt"

	log "github.com/sirupsen/logrus"
)

type IMarketPay interface {
	CreateOrder(name string, value string) (*Order, error)
	GetOrder(hash string) (*Order, error)
	IsPaid(order Order) bool
}

type MarketPay struct {
	endpoint  string
	publicKey string
}

func New() *MarketPay {
	sandboxURL := "https://api.truust.io/"
	sandboxPublicKey := "pk_production_Q2F2VlMxSEk="

	return &MarketPay{
		sandboxURL,
		sandboxPublicKey,
	}
}

func NewSandbox() *MarketPay {
	sandboxURL := "https://api-sandbox.truust.io/1.0"
	sandboxPublicKey := "pk_stage_ZkNpNElWeEg="

	return &MarketPay{
		sandboxURL,
		sandboxPublicKey,
	}
}

func (b *MarketPay) CreateOrder(
	name string,
	value string,
) (*Order, error) {
	log.Infof("[Marketpay] Create order name %v value %v", name, value)

	url := b.endpoint + "/express"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("name", name)
	_ = writer.WriteField("amount", value)
	_ = writer.WriteField("source", b.publicKey)
	_ = writer.WriteField("trustee_confirmed_url", "https://www.goalrev.com/purchasesuccess")
	_ = writer.WriteField("trustee_denied_url", "https://www.goalrev.com/purchasefailure")
	_ = writer.WriteField("settlor_confirmed_url", "https://www.goalrev.com/sellsuccess")
	_ = writer.WriteField("settlor_denied_url", "https://www.goalrev.com/sellfailure")
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
		return nil, err
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
