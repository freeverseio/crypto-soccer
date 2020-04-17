package consumer

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func CreateAuction(tx *sql.Tx, in input.CreateAuctionInput) error {
	signerAddress, err := in.SignerAddress()
	if err != nil {
		return err
	}
	auction := storage.NewAuction()
	id, err := in.ID()
	if err != nil {
		return err
	}
	auction.ID = string(id)
	auction.Rnd = int(in.Rnd)
	auction.PlayerID = in.PlayerId
	auction.CurrencyID = int(in.CurrencyId)
	auction.Price = int(in.Price)
	if auction.ValidUntil, err = strconv.ParseInt(in.ValidUntil, 10, 64); err != nil {
		fmt.Printf("%d of type %T", auction.ValidUntil, auction.ValidUntil)
	}
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
