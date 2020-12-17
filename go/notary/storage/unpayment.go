package storage

type Unpayment struct {
	Owner               string
	LastTimeOfUnpayment string
	Notified            bool
}

func NewUnpayment() *Unpayment {
	unpayment := Unpayment{}
	return &unpayment
}
