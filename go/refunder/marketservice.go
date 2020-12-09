package refunder

type MarketService interface {
	Begin() (MarketTx, error)
}

type MarketTx interface {
}
