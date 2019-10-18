package storage

type BidState string

const (
	BID_FILED   BidState = "FILED"
	BID_PAYING  BidState = "PAYING"
	BID_PAID    BidState = "PAID"
	BID_EXPIRED BidState = "EXPIRED"
)
