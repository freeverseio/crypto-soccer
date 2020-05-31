package consumer

import (
	"database/sql"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func CreateBid(tx *sql.Tx, in input.CreateBidInput) error {
	service := postgres.NewAuctionService(tx)
	auction, err := service.Auction(string(in.AuctionId))
	if err != nil {
		return err
	}

	if auction == nil {
		return fmt.Errorf("No auction for bid %v", in)
	}

	bid := storage.NewBid()
	bid.AuctionID = string(in.AuctionId)
	bid.ExtraPrice = int64(in.ExtraPrice)
	bid.Rnd = int64(in.Rnd)
	bid.TeamID = in.TeamId
	bid.Signature = in.Signature
	bid.State = storage.BidAccepted
	bid.StateExtra = ""
	bid.PaymentID = ""
	bid.PaymentURL = ""
	bid.PaymentDeadline = 0

	return service.Bid().Insert(*bid)
}
