package purchasevoider

type MarketService interface {
	GetPlayerIdByPurchaseToken(token string) (string, error)
}
