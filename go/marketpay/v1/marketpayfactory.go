package v1

type IMarketPay interface {
	CreateOrder(name string, value string) (*Order, error)
	GetOrder(hash string) (*Order, error)
	IsPaid(order Order) bool
}

type IMarketPayContext interface {
	GetEndPoint() string
	GetPublicKey() string
}

type MarketPayFactory struct{}

func (factory MarketPayFactory) Create(context IMarketPayContext) (IMarketPay, error) {
	switch v := context.(type) {
	case mockMarketPayContext:
		return NewMockMarketPay(v)
	}
	return NewMarketPay(context.(MarketPayContext))
}
