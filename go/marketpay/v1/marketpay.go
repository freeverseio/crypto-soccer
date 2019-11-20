package v1

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"fmt"
)

type Order struct {
	Name                string      `json:"name"`
	Amount              string      `json:"amount"`
	Currency            string      `json:"currency"`
	WithShipping        int         `json:"with_shipping"`
	WithInsurance       int         `json:"with_insurance"`
	ConditionsTrustee   string      `json:"conditions_trustee"`
	ConditionsSettlor   string      `json:"conditions_settlor"`
	ConditionsCredit    string      `json:"conditions_credit"`
	ConditionsInsurance string      `json:"conditions_insurance"`
	Status              string      `json:"status"`
	Tag                 interface{} `json:"tag"`
	ID                  int         `json:"id"`
	FormatAmount        string      `json:"format_amount"`
	FormatFee           string      `json:"format_fee"`
	FormatFeeShipping   string      `json:"format_fee_shipping"`
	FormatFeeInsurance  string      `json:"format_fee_insurance"`
	FormatPayinAmount   string      `json:"format_payin_amount"`
	FormatCreditAmount  string      `json:"format_credit_amount"`
	FormatPayoutAmount  string      `json:"format_payout_amount"`
	FormatCreatedAt     string      `json:"format_created_at"`
	StatusNicename      string      `json:"status_nicename"`
	TrusteeShortlink    struct {
		Hash     string `json:"hash"`
		ShortURL string `json:"short_url"`
	} `json:"trustee_shortlink"`
	SettlorShortlink struct {
		Hash     string `json:"hash"`
		ShortURL string `json:"short_url"`
	} `json:"settlor_shortlink"`
	Symbol    string `json:"symbol"`
	Shortlink struct {
		TrusteeConfirmedURL string      `json:"trustee_confirmed_url"`
		TrusteeDeniedURL    string      `json:"trustee_denied_url"`
		SettlorConfirmedURL string      `json:"settlor_confirmed_url"`
		SettlorDeniedURL    string      `json:"settlor_denied_url"`
		ID                  int         `json:"id"`
		Hash                string      `json:"hash"`
		Visits              interface{} `json:"visits"`
		ShortURL            string      `json:"short_url"`
		Source              struct {
			ReferenceID        string      `json:"reference_id"`
			CountryCode        string      `json:"country_code"`
			Name               string      `json:"name"`
			Slug               string      `json:"slug"`
			FeeShipping1       int         `json:"fee_shipping_1"`
			FeeShipping2       int         `json:"fee_shipping_2"`
			FeeShipping3       int         `json:"fee_shipping_3"`
			DefaultParcelSize  string      `json:"default_parcel_size"`
			PrimaryColor       interface{} `json:"primary_color"`
			SecondaryColor     interface{} `json:"secondary_color"`
			TemplateFolder     string      `json:"template_folder"`
			AllowInsurance     int         `json:"allow_insurance"`
			AllowShipping      int         `json:"allow_shipping"`
			AllowSelectParcel  int         `json:"allow_select_parcel"`
			AllowCredit        int         `json:"allow_credit"`
			PhotoURL           string      `json:"photo_url"`
			CurrentBillingPlan interface{} `json:"current_billing_plan"`
			PublicKey          string      `json:"public_key"`
			IsActive           bool        `json:"is_active"`
			TaxRate            int         `json:"tax_rate"`
		} `json:"source"`
	} `json:"shortlink"`
	Image interface{} `json:"image"`
}

const sandboxURL = "https://api-sandbox.truust.io/1.0"
const sandboxPublicKey = "pk_stage_ZkNpNElWeEg="

type MarketPay struct {
	endpoint  string
	publicKey string
}

func New() (*MarketPay, error) {
	return &MarketPay{
		sandboxURL,
		sandboxPublicKey,
	}, nil
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
	_ = writer.WriteField("trustee_confirmed_url", "https://ninjaflex.com/")
	_ = writer.WriteField("trustee_denied_url", "https://ninjaflex.com/")
	_ = writer.WriteField("settlor_confirmed_url", "https://ninjaflex.com/")
	_ = writer.WriteField("settlor_denied_url", "https://ninjaflex.com/")
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
