package auctionmachine

type AssetFrozen struct {
}

func NewAssetFrozen() State {
	return &AssetFrozen{}
}

func (b *AssetFrozen) Process(m *AuctionMachine) error {
	return nil
}
