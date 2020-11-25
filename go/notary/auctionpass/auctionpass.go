package auctionpass

type AuctionPass struct {
	owner              string
	purchasedForTeamId string
	productId          string
	ack                bool
}

func NewAuctionPass(
	owner string,
	purchasedForTeamId string,
	productId string,
	ack bool,
) *AuctionPass {
	auctionPass := AuctionPass{}
	auctionPass.owner = owner
	auctionPass.purchasedForTeamId = purchasedForTeamId
	auctionPass.productId = productId
	auctionPass.ack = ack
	return &auctionPass
}

func (b AuctionPass) Owner() string {
	return b.owner
}

func (b AuctionPass) PurchasedForTeamId() string {
	return b.purchasedForTeamId
}

func (b AuctionPass) ProductId() string {
	return b.productId
}

func (b AuctionPass) Ack() bool {
	return b.ack
}
