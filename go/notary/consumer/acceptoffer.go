package consumer

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func AcceptOffer(service storage.StorageService, tx *sql.Tx, in input.AcceptOfferInput) error {
	offerID, err := strconv.ParseInt(string(in.OfferId), 10, 64)
	if err != nil {
		return err
	}

	offer, err := service.Offer(tx, offerID)
	if err != nil {
		return err
	}

	if offer != nil && offer.State != storage.OfferStarted {
		return errors.New("Auctions can only be created for offers in Started state")
	}

	if offer != nil && offer.ValidUntil < time.Now().Unix() {
		offer.State = storage.OfferEnded
		offer.StateExtra = "Offer expired when accepting"
		if err = service.OfferUpdate(tx, *offer); err != nil {
			return err
		}
		return errors.New("Associated Offer is expired")
	}

	auction := storage.NewAuction()
	id, err := in.AuctionID()
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
	if err = service.AuctionInsert(tx, *auction); err != nil {
		return err
	}

	if offer != nil && offer.ID != 0 {
		offer.State = storage.OfferAccepted
		offer.StateExtra = ""
		offer.AuctionID = auction.ID
		offer.Seller = auction.Seller

		if err = service.OfferUpdate(tx, *offer); err != nil {
			return err
		}

		bid := input.CreateBidInput{
			Signature:  offer.Signature,
			AuctionId:  graphql.ID(offer.AuctionID),
			ExtraPrice: 0,
			Rnd:        int32(offer.Rnd),
			TeamId:     offer.TeamID,
		}

		err = CreateBid(service, tx, bid)

		if err != nil {
			log.Error(err)
			offer.State = storage.OfferFailed
			offer.StateExtra = "Could not create bid"
			service.OfferUpdate(tx, *offer)
			return err
		}
	}

	return nil
}
