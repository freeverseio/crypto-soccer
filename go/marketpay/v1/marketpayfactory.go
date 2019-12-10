package v1

type IMarketPay interface {
	CreateOrder(name string, value string) (*Order, error)
	GetOrder(hash string) (*Order, error)
	IsPaid(order Order) bool
}

type MarketPayFactory struct{}

func (factory MarketPayFactory) Create(endPoint string, publicKey string) (IMarketPay, error) {
	if len(endPoint) == 0 || len(publicKey) == 0 {
		return &mockMarketPay{}, nil
	}
	return NewMarketPay(endPoint, publicKey)
}
