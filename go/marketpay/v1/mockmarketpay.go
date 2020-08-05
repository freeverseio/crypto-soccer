package v1

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/freeverseio/crypto-soccer/go/marketpay"
)

type MockMarketPay struct {
	orders      map[string]*marketpay.Order
	orderStatus OrderStatus
}

func NewMockMarketPay() *MockMarketPay {
	return &MockMarketPay{
		orders:      make(map[string]*marketpay.Order),
		orderStatus: DRAFT,
	}
}

func (m *MockMarketPay) SetOrderStatus(s OrderStatus) {
	m.orderStatus = s
	for k, v := range m.orders {
		v.Status = m.orderStatus.String()
		m.orders[k] = v
	}
}

func (m *MockMarketPay) CreateOrder(name string, value string) (*marketpay.Order, error) {
	order := &marketpay.Order{}
	order.Status = m.orderStatus.String()
	order.Name = name
	order.Amount = value
	hasher := sha256.New()
	hasher.Write([]byte(name))
	hash := hex.EncodeToString(hasher.Sum(nil))[:6]
	order.TrusteeShortlink.Hash = hash
	order.TrusteeShortlink.ShortURL = "https://trustee.io/" + hash
	order.SettlorShortlink.ShortURL = "https://settlor.io/" + hash
	m.orders[hash] = order
	return order, nil
}

func (m *MockMarketPay) GetOrder(hash string) (*marketpay.Order, error) {
	order, ok := m.orders[hash]
	if !ok {
		return nil, errors.New("Could not find order")
	}
	return order, nil
}

func (m *MockMarketPay) IsPaid(order marketpay.Order) bool {
	return order.Status == PUBLISHED.String()
}

func (b *MockMarketPay) ValidateOrder(hash string) (string, error) {
	return "not implemented", nil
}
