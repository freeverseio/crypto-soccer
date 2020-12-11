package storage

type AuctionPassPlaystoreOrderState string

const (
	AuctionPassPlaystoreOrderOpen         AuctionPassPlaystoreOrderState = "open"
	AuctionPassPlaystoreOrderAcknowledged AuctionPassPlaystoreOrderState = "acknowledged"
	AuctionPassPlaystoreOrderComplete     AuctionPassPlaystoreOrderState = "complete"
	AuctionPassPlaystoreOrderRefunding    AuctionPassPlaystoreOrderState = "refunding"
	AuctionPassPlaystoreOrderRefunded     AuctionPassPlaystoreOrderState = "refunded"
	AuctionPassPlaystoreOrderFailed       AuctionPassPlaystoreOrderState = "failed"
)

type AuctionPassPlaystoreOrder struct {
	OrderId       string
	PackageName   string
	ProductId     string
	PurchaseToken string
	TeamId        string
	Owner         string
	Signature     string
	State         AuctionPassPlaystoreOrderState
	StateExtra    string
}

func NewAuctionPassPlaystoreOrder() *AuctionPassPlaystoreOrder {
	order := AuctionPassPlaystoreOrder{}
	order.State = AuctionPassPlaystoreOrderOpen
	return &order
}
