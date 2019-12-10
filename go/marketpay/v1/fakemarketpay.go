package v1

type fakeMarketPay struct {
}

func (b *fakeMarketPay) CreateOrder(name string, value string) (*Order, error) {
	order := &Order{}
	order.Status = "DRAFT"
	order.Name = name
	order.Amount = value
	return order, nil
}

func (b *fakeMarketPay) GetOrder(hash string) (*Order, error) {
	order := &Order{}
	order.Status = "DRAFT"
	return order, nil
}

func (b *fakeMarketPay) IsPaid(order Order) bool {
	return false
}
