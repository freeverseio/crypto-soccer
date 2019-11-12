package marketpay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type Customer struct {
	Data struct {
		ID          int         `json:"id"`
		Self        string      `json:"self"`
		Name        interface{} `json:"name"`
		Email       interface{} `json:"email"`
		Prefix      string      `json:"prefix"`
		Phone       string      `json:"phone"`
		Tag         interface{} `json:"tag"`
		Metadata    interface{} `json:"metadata"`
		CreatedAt   time.Time   `json:"created_at"`
		Connections struct {
			Wallet       string `json:"wallet"`
			Bankaccounts string `json:"bankaccounts"`
		} `json:"connections"`
	} `json:"data"`
}
type Order struct {
	Data struct {
		ID             int           `json:"id"`
		Self           string        `json:"self"`
		PublicID       string        `json:"public_id"`
		Name           string        `json:"name"`
		Value          string        `json:"value"`
		Currency       string        `json:"currency"`
		Amount         string        `json:"amount"`
		PayinAmount    string        `json:"payin_amount"`
		PayoutAmount   string        `json:"payout_amount"`
		FeesAmount     string        `json:"fees_amount"`
		BuyerLink      string        `json:"buyer_link"`
		SellerLink     string        `json:"seller_link"`
		Status         string        `json:"status"`
		Sandbox        int           `json:"sandbox"`
		StatusNicename string        `json:"status_nicename"`
		Images         []interface{} `json:"images"`
		Tag            interface{}   `json:"tag"`
		Metadata       interface{}   `json:"metadata"`
		Refund         struct {
			Status      interface{} `json:"status"`
			ReferenceID interface{} `json:"reference_id"`
		} `json:"refund"`
		CreatedAt   time.Time   `json:"created_at"`
		PublishedAt interface{} `json:"published_at"`
		AcceptedAt  interface{} `json:"accepted_at"`
		ValidatedAt interface{} `json:"validated_at"`
		ReleasedAt  interface{} `json:"released_at"`
		Connections struct {
			Wallet string `json:"wallet"`
			Buyer  string `json:"buyer"`
			Payins string `json:"payins"`
			Seller string `json:"seller"`
		} `json:"connections"`
	} `json:"data"`
}

const sandboxURL = "https://api-sandbox.truust.io"
const sandBoxBearerToken = "Bearer sk_stage_NCzkqJwQTNVStxDxVxmSflVv"

const productionURL = "https://api.truust.io"
const productionBearerToken = "Bearer sk_production_AjWbpOPwS3HNi821Ma9mIgA2"

type MarketPay struct {
	endpoint    string
	bearerToken string
}

func New() (*MarketPay, error) {
	return &MarketPay{
		sandboxURL,
		sandBoxBearerToken,
	}, nil
}

func (b *MarketPay) CreateCustomer(
	prefix string,
	phoneNumber string,
) (*Customer, error) {
	url := b.endpoint + "/2.0/customers"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("prefix", prefix)
	_ = writer.WriteField("phone", phoneNumber)
	err := writer.Close()
	if err != nil {
		return nil, err
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

	customer := &Customer{}
	err = json.Unmarshal(body, customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (b *MarketPay) CreateOrder(
	seller *Customer,
	buyer *Customer,
	name string,
	value string,
) (*Order, error) {
	url := b.endpoint + "/2.0/orders"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("buyer_id", fmt.Sprintf("%d", buyer.Data.ID))
	_ = writer.WriteField("seller_id", fmt.Sprintf("%d", seller.Data.ID))
	_ = writer.WriteField("name", name)
	_ = writer.WriteField("value", value)
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

	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	// fmt.Println(string(body))
	order := &Order{}
	err = json.Unmarshal(body, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (b *MarketPay) GetOrder(orderID int) (*Order, error) {
	url := b.endpoint + "/2.0/orders/" + strconv.Itoa(orderID)
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
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	// fmt.Println(string(body))
	order := &Order{}
	err = json.Unmarshal(body, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (b *MarketPay) IsPaid(order *Order) bool {
	return order.Data.Status == "PUBLISHED"
}
