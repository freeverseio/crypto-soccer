package v1

import (
	"errors"
)

type mockMarketPay struct {
	orders map[string]*Order
}

func NewMockMarketPay() (*mockMarketPay, error) {
	return &mockMarketPay{
		orders: make(map[string]*Order),
	}, nil
}

func (b *mockMarketPay) CreateOrder(name string, value string) (*Order, error) {
	order := &Order{}
	order.Status = "DRAFT"
	order.Name = name
	order.Amount = value
	order.TrusteeShortlink.Hash = name
	b.orders[name] = order
	return order, nil
}

func (b *mockMarketPay) GetOrder(hash string) (*Order, error) {
	order, ok := b.orders[hash]
	if !ok {
		return nil, errors.New("Could not find order")
	}
	return order, nil
}

func (b *mockMarketPay) IsPaid(order Order) bool {
	return false
}
