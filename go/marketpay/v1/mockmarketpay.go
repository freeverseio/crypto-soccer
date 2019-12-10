package v1

type mockMarketPay struct {
}

func (b *mockMarketPay) CreateOrder(name string, value string) (*Order, error) {
	order := &Order{}
	order.Status = "DRAFT"
	order.Name = name
	order.Amount = value
	return order, nil
}

func (b *mockMarketPay) GetOrder(hash string) (*Order, error) {
	order := &Order{}
	order.Status = "DRAFT"
	return order, nil
}

func (b *mockMarketPay) IsPaid(order Order) bool {
	return false
}
