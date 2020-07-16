package storage

type OfferState string

const (
	OfferStarted   OfferState = "started"
	OfferEnded     OfferState = "ended"
	OfferCancelled OfferState = "cancelled"
	OfferFailed    OfferState = "failed"
	OfferAccepted  OfferState = "accepted"
)

type Offer struct {
	ID         string
	PlayerID   string
	CurrencyID int
	Price      int64
	Rnd        int64
	ValidUntil int64
	Signature  string
	State      OfferState
	StateExtra string
	Seller     string
	Buyer      string
	AuctionID  string
	TeamID     string
}

func NewOffer() *Offer {
	offer := Offer{}
	offer.State = OfferStarted
	return &offer
}

type OfferService interface {
	PendingOffers() ([]Offer, error)
	Offer(ID string) (*Offer, error)
	Insert(offer Offer) error
	Update(offer Offer) error
}
