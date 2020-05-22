package consumer

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
)

func CancelAuction(tx *sql.Tx, in input.CancelAuctionInput) error {
	service := postgres.NewAuctionService(tx)
	auction, err := service.Auction(string(in.AuctionId))
	if err != nil {
		return err
	}
	if auction == nil {
		return errors.New("trying to cancel a nil auction")
	}
	if auction.State != storage.AuctionStarted {
		return fmt.Errorf("not possible to cancel an auction in state %v", auction.State)
	}

	auction.State = storage.AuctionCancelled
	return service.Update(*auction)
}
