package consumer

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func CreateOffer(tx *sql.Tx, in input.CreateOfferInput, contracts contracts.Contracts) error {
	offer := storage.NewOffer()
	id, err := in.ID(contracts)
	if err != nil {
		return err
	}
	offer.ID = string(id)
	offer.Rnd = int64(in.Rnd)
	offer.PlayerID = in.PlayerId
	offer.CurrencyID = int(in.CurrencyId)
	offer.Price = int64(in.Price)
	if offer.ValidUntil, err = strconv.ParseInt(in.ValidUntil, 10, 64); err != nil {
		fmt.Printf("%d of type %T", offer.ValidUntil, offer.ValidUntil)
	}
	offer.Signature = in.Signature
	offer.State = storage.OfferStarted
	offer.StateExtra = ""
	signerAddress, err := in.SignerAddress(contracts)
	if err != nil {
		return err
	}
	offer.Buyer = signerAddress.Hex()
	offer.Seller = in.Seller
	service := postgres.NewOfferHistoryService(tx)
	if err = service.Insert(*offer); err != nil {
		return err
	}

	return nil
}
