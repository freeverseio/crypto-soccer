package consumer

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func CreateAuction(tx *sql.Tx, in input.CreateAuctionInput) error {
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

	if in.OfferId != "" {
		offer := storage.NewOffer()
		offer.ID = in.OfferId
		offer.State = storage.OfferAccepted
		offer.StateExtra = ""
		offer.AuctionID = auction.ID
		offer.Seller = auction.Seller

		service := postgres.NewOfferService(tx)
		if err = service.Update(*offer); err != nil {
			return err
		}
	}

	return nil
}
