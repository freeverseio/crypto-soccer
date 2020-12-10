package purchasevoider

type MarketService interface {
	Begin() (MarketTx, error)
}

type MarketTx interface {
}
