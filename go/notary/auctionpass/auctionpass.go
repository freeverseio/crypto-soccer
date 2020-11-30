package auctionpass

type AuctionPass struct {
	owner              string
	purchasedForTeamId string
	productId          string
	ack                bool
}

func NewAuctionPass(
	owner string,
) *AuctionPass {
	auctionPass := AuctionPass{}
	auctionPass.owner = owner
	return &auctionPass
}

func (b AuctionPass) Owner() string {
	return b.owner
}
