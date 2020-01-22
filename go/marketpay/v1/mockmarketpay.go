package v1

import (
	"errors"
)

type MockMarketPay struct {
	orders      map[string]*Order
	orderStatus OrderStatus
}

func NewMockMarketPay() *MockMarketPay {
	return &MockMarketPay{
		orders:      make(map[string]*Order),
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

func (m *MockMarketPay) CreateOrder(name string, value string) (*Order, error) {
	order := &Order{}
	order.Status = m.orderStatus.String()
	order.Name = name
	order.Amount = value
	order.TrusteeShortlink.Hash = name
	order.TrusteeShortlink.ShortURL = "https://trustee.io/" + name
	order.SettlorShortlink.ShortURL = "https://settlor.io/" + name
	m.orders[name] = order
	return order, nil
}

func (m *MockMarketPay) GetOrder(hash string) (*Order, error) {
	order, ok := m.orders[hash]
	if !ok {
		return nil, errors.New("Could not find order")
	}
	return order, nil
}

func (m *MockMarketPay) IsPaid(order Order) bool {
	return order.Status == PUBLISHED.String()
}
