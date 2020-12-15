package storage

type Unpayment struct {
	Owner               string
	NumOfUnpayments     int
	LastTimeOfUnpayment string
}

func NewUnpayment() *Unpayment {
	unpayment := Unpayment{}
	return &unpayment
}
