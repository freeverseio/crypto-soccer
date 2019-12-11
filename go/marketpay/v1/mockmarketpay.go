package v1

import (
	"errors"
)

type mockMarketPay struct {
	orders  map[string]*Order
	context mockMarketPayContext
}
type mockMarketPayContext struct {
	states []OrderStatus
	idx    int
}

func NewMockMarketPayContext(states []OrderStatus) mockMarketPayContext {
	return mockMarketPayContext{states, 0}
}

func (c mockMarketPayContext) GetEndPoint() string {
	return ""
}

func (c mockMarketPayContext) GetPublicKey() string {
	return ""
}

func (c *mockMarketPayContext) NextOrderStatus() OrderStatus {
	state := c.states[c.idx]
	if c.idx < len(c.states) {
		c.idx++
	}
	return state
}

func NewMockMarketPay(context mockMarketPayContext) (*mockMarketPay, error) {
	return &mockMarketPay{
		orders:  make(map[string]*Order),
		context: context,
	}, nil
}

func (b *mockMarketPay) CreateOrder(name string, value string) (*Order, error) {
	order := &Order{}
	order.Status = DRAFT.String()
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
	order.Status = b.context.NextOrderStatus().String()
	return order, nil
}

func (b *mockMarketPay) IsPaid(order Order) bool {
	return order.Status == PUBLISHED.String()
}
