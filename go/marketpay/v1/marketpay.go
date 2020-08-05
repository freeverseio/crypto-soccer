package v1

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"fmt"

	"github.com/freeverseio/crypto-soccer/go/marketpay"
	log "github.com/sirupsen/logrus"
)

type MarketPay struct {
	endpoint    string
	endpoint2   string
	publicKey   string
	bearerToken string
}

func New(pk string) *MarketPay {
	market := MarketPay{}
	market.endpoint = "https://api.truust.io/1.0"
	market.endpoint2 = "https://api.truust.io/2.0"
	market.publicKey = "pk_production_Q2F2VlMxSEk="
	market.bearerToken = "Bearer " + pk
	return &market
}

func NewSandbox() *MarketPay {
	market := MarketPay{}
	market.endpoint = "https://api-sandbox.truust.io/1.0"
	market.endpoint2 = "https://api-sandbox.truust.io/2.0"
	market.publicKey = "pk_stage_ZkNpNElWeEg="
	market.bearerToken = "Bearer sk_stage_NCzkqJwQTNVStxDxVxmSflVv"
	return &market
}

func (b *MarketPay) CreateOrder(
	name string,
	value string,
) (*marketpay.Order, error) {
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

	order := &marketpay.Order{}
	err = json.Unmarshal(body, order)
	if err != nil {
		return nil, err
	}

	// log.Infof("[marketpayV1] order %+v created", order)
	return order, nil
}

func (b *MarketPay) GetOrder(hash string) (*marketpay.Order, error) {
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

	order := &marketpay.Order{}
	err = json.Unmarshal(body, order)
	if err != nil {
		log.Error(order)
		return nil, err
	}
	return order, nil
}

func (b *MarketPay) IsPaid(order marketpay.Order) bool {
	return order.Status == "PUBLISHED"
}

func (b *MarketPay) ValidateOrder(hash string) (string, error) {
	order, err := b.GetOrder(hash)
	if err != nil {
		return "", err
	}
	log.Infof("Validate order ID %v", order.ID)
	url := b.endpoint2 + "/orders/" + order.ID + "/validate"
	method := "POST"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", b.bearerToken)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// log.Info(string(body))

	// order := &Order{}
	// err = json.Unmarshal(body, order)
	// if err != nil {
	// 	log.Error(order)
	// 	return "", err
	// }
	return string(body), nil
}
