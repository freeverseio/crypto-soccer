package storage

type MarketService interface {
	Begin() (MarketTx, error)
}

type MarketTx interface {
}
