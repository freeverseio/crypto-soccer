package consumer

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func CreateAuction(tx *sql.Tx, in input.CreateAuctionInput) error {

	offerService := postgres.NewOfferService(tx)

	offer, err := offerService.OfferByRndPrice(in.Rnd, in.Price)
	if err != nil {
		return err
	}
	if offer != nil && offer.State != storage.OfferStarted {
		return errors.New("Auctions can only be created for offers in Started state")
	}
	if offer != nil && offer.ValidUntil < time.Now().Unix() {
		offer.State = storage.OfferEnded
		offer.StateExtra = "Offer expired when accepting"
		if err = offerService.Update(*offer); err != nil {
			return err
		}
		return errors.New("Associated Offer is expired")
	}

	auction := storage.NewAuction()
	id, err := in.ID()
	if err != nil {
		return err
	}
	auction.ID = string(id)
	auction.Rnd = int64(in.Rnd)
	auction.PlayerID = in.PlayerId
	auction.CurrencyID = int(in.CurrencyId)
	auction.Price = int64(in.Price)
	if auction.ValidUntil, err = strconv.ParseInt(in.ValidUntil, 10, 64); err != nil {
		return fmt.Errorf("invalid validUntil %v", in.ValidUntil)
	}
	auction.Signature = in.Signature
	auction.State = storage.AuctionStarted
	auction.StateExtra = ""
	auction.PaymentURL = ""
	signerAddress, err := in.SignerAddress()
	if err != nil {
		return err
	}
	auction.Seller = signerAddress.Hex()
	service := postgres.NewAuctionHistoryService(tx)
	if err = service.Insert(*auction); err != nil {
		return err
	}

	// offer was necessarily in OfferStarted state, due to check above
	if offer != nil && offer.ID != "" {
		offer.State = storage.OfferAccepted
		offer.StateExtra = ""
		offer.AuctionID = auction.ID
		offer.Seller = auction.Seller

		if err = offerService.Update(*offer); err != nil {
			return err
		}
	}

	return nil
}
