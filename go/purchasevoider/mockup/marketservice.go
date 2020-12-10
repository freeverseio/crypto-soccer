package mockup

type MarketService struct {
	GetPlayerIdByPurchaseTokenFunc func(token string) (string, error)
}

func (b MarketService) GetPlayerIdByPurchaseToken(token string) (string, error) {
	return b.GetPlayerIdByPurchaseTokenFunc(token)
}
