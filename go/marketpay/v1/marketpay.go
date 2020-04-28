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
	endpoint    string
	publicKey   string
	bearerToken string
}

func New() *MarketPay {
	market := MarketPay{}
	market.endpoint = "https://api.truust.io/1.0"
	market.publicKey = "pk_production_Q2F2VlMxSEk="
	market.bearerToken = "Bearer sk_production_AjWbpOPwS3HNi821Ma9mIgA2"

	return &market
}

func NewSandbox() *MarketPay {
	market := MarketPay{}
	market.endpoint = "https://api-sandbox.truust.io/1.0"
	market.publicKey = "pk_stage_ZkNpNElWeEg="
	market.bearerToken = "Bearer sk_stage_NCzkqJwQTNVStxDxVxmSflVv"

	return &market
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

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", b.bearerToken)
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

	order := &Order{}
	err = json.Unmarshal(body, order)
	if err != nil {
		return nil, err
	}

	// log.Infof("[marketpayV1] order %+v created", order)
	return order, nil
}

func (b *MarketPay) GetOrder(hash string) (*Order, error) {
	log.Infof("Getting order %v", hash)
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
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", b.bearerToken)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	order := &Order{}
	err = json.Unmarshal(body, order)
	if err != nil {
		log.Error(order)
		return nil, err
	}
	return order, nil
}

func (b *MarketPay) IsPaid(order Order) bool {
	return order.Status == "PUBLISHED"
}
