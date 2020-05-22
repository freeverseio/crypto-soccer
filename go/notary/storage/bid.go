package storage

type BidState string

const (
	BIDREFUSED  BidState = "REFUSED"
	BidAccepted BidState = "accepted"
	BidPaid     BidState = "paid"
	BidPaying   BidState = "paying"
	BidFailed   BidState = "failed"
)

type Bid struct {
	AuctionID       string
	ExtraPrice      int64
	Rnd             int64
	TeamID          string
	Signature       string
	State           BidState
	StateExtra      string
	PaymentID       string
	PaymentURL      string
	PaymentDeadline int64
}

type BidService interface {
	Bids(auctionId string) ([]Bid, error)
	Insert(bid Bid) error
	Update(bid Bid) error
}

func NewBid() *Bid {
	bid := Bid{}
	bid.State = BidAccepted
	return &bid
}

func FindBids(bids []Bid, state BidState) []*Bid {
	result := []*Bid{}
	for i := range bids {
		if bids[i].State == state {
			result = append(result, &bids[i])
		}
	}
	return result
}
