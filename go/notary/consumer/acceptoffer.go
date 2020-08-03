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

func highestOffer(offers []storage.Offer) (*storage.Offer, error) {
	length := len(offers)
	if length == 0 {
		return nil, errors.New("There are no offers for this playerId")
	}
	if length == 1 {
		return &offers[0], nil
	}

	idx := -1
	price := int64(-1)
	for i, offer := range offers {
		if offer.State == storage.OfferStarted {
			if idx == -1 {
				idx = i
				price = offer.Price
			} else {
				if offer.Price > price {
					idx = i
					price = offer.Price
				}
			}
		}
	}
	if idx == -1 {
		return nil, errors.New("There are not acceptable offers")
	}

	return &offers[idx], nil
}

func AcceptOffer(service storage.StorageService, tx *sql.Tx, in input.AcceptOfferInput) error {
	offers, err := service.OffersByPlayerId(string(in.PlayerId))

	highestOffer, err := highestOffer(offers)
	if err != nil {
		return err
	}

	if highestOffer.ID != string(in.OfferId) {
		return errors.New("You can only accept highest offer")
	}

	offer := highestOffer
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

	if offer != nil && offer.ID != "" {
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
