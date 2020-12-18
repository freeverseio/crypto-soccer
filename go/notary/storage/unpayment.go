package storage

type Unpayment struct {
	Id              int64
	Owner           string
	TimeOfUnpayment string
	AuctionId       string
	Notified        bool
}

func NewUnpayment() *Unpayment {
	unpayment := Unpayment{}
	return &unpayment
}
