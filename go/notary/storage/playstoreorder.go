package storage

type PlaystoreOrderState string

const (
	PlaystoreOrderOpen         PlaystoreOrderState = "open"
	PlaystoreOrderAcknowledged PlaystoreOrderState = "acknowledged"
	PlaystoreOrderComplete     PlaystoreOrderState = "complete"
	PlaystoreOrderRefunding    PlaystoreOrderState = "refunding"
	PlaystoreOrderRefunded     PlaystoreOrderState = "refunded"
	PlaystoreOrderFailed       PlaystoreOrderState = "failed"
)

type PlaystoreOrder struct {
	OrderId       string
	PackageName   string
	ProductId     string
	PurchaseToken string
	PlayerId      string
	TeamId        string
	Signature     string
	State         PlaystoreOrderState
	StateExtra    string
}

func NewPlaystoreOrder() *PlaystoreOrder {
	order := PlaystoreOrder{}
	order.State = PlaystoreOrderOpen
	return &order
}
