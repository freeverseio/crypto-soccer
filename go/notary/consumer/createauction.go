package consumer

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func CreateAuction(tx *sql.Tx, in input.CreateAuctionInput) error {
	signerAddress, err := in.SignerAddress()
	if err != nil {
		return err
	}
	auction := storage.NewAuction()
	auction.ID = in.ID()
	auction.Rnd = int(in.Rnd)
	auction.PlayerID = in.PlayerId
	auction.CurrencyID = int(in.CurrencyId)
	auction.Price = int(in.Price)
	auction.ValidUntil = in.ValidUntil
	auction.Signature = in.Signature
	auction.State = storage.AuctionStarted
	auction.StateExtra = ""
	auction.PaymentURL = ""
	auction.Seller = signerAddress.Hex()
	if err = auction.Insert(tx); err != nil {
		return err
	}

	return nil
}
