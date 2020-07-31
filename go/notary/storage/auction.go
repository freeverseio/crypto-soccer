package storage

type AuctionState string

const (
	AuctionAssetFrozen        AuctionState = "asset_frozen"
	AuctionPaying             AuctionState = "paying"
	AuctionWithdrableBySeller AuctionState = "withadrable_by_seller"
	AuctionWithdrableByBuyer  AuctionState = "withadrable_by_buyer"
	AuctionStarted            AuctionState = "started"
	AuctionEnded              AuctionState = "ended"
	AuctionCancelled          AuctionState = "cancelled"
	AuctionFailed             AuctionState = "failed"
	AuctionValidation         AuctionState = "validation"
)

type Auction struct {
	ID         string
	PlayerID   string
	CurrencyID int
	Price      int64
	Rnd        int64
	ValidUntil int64
	Signature  string
	State      AuctionState
	StateExtra string
	PaymentURL string
	Seller     string
}

func NewAuction() *Auction {
	auction := Auction{}
	auction.State = AuctionStarted
	return &auction
}

type AuctionService interface {
	PendingAuctions() ([]Auction, error)
	Auction(ID string) (*Auction, error)
	AuctionsByPlayerId(ID string) ([]Auction, error)
	Insert(auction Auction) error
	Update(auction Auction) error
	Bid() BidService
}
