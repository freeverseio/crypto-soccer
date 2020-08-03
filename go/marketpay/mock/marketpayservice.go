package mock

import "github.com/freeverseio/crypto-soccer/go/marketpay"

type MarketPayMock struct {
	CreateOrderFunc   func(name string, value string) (*marketpay.Order, error)
	GetOrderFunc      func(hash string) (*marketpay.Order, error)
	IsPaidFunc        func(order marketpay.Order) bool
	ValidateOrderFunc func(hash string) (string, error)
}

func (b *MarketPayMock) CreateOrder(name string, value string) (*marketpay.Order, error) {
	return b.CreateOrderFunc(name, value)
}

func (b *MarketPayMock) GetOrder(hash string) (*marketpay.Order, error) {
	return b.GetOrderFunc(hash)
}

func (b *MarketPayMock) IsPaid(order marketpay.Order) bool {
	return b.IsPaidFunc(order)
}

func (b *MarketPayMock) ValidateOrder(hash string) (string, error) {
	return b.ValidateOrderFunc(hash)
}
