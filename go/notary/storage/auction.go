package storage

import (
	"database/sql"
	"math/big"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type AuctionState string

const (
	AUCTION_STARTED             AuctionState = "STARTED"
	AUCTION_ASSET_FROZEN        AuctionState = "ASSET_FROZEN"
	AUCTION_PAYING              AuctionState = "PAYING"
	AUCTION_PAID                AuctionState = "PAID"
	AUCTION_NO_BIDS             AuctionState = "NO_BIDS"
	AUCTION_CANCELLED_BY_SELLER AuctionState = "CANCELLED_BY_SELLER"
	AUCTION_WITHDRAWAL          AuctionState = "WITHDRAWAL"
	AUCTION_FAILED              AuctionState = "FAILED"
	AuctionEnded                AuctionState = "Ended"
)

type Auction struct {
	UUID       uuid.UUID
	PlayerID   *big.Int
	CurrencyID uint8
	Price      *big.Int
	Rnd        *big.Int
	ValidUntil *big.Int
	Signature  string
	State      AuctionState
	StateExtra string
	PaymentURL string
	Seller     string
}

func (b *Storage) CreateAuction(order Auction) error {
	log.Infof("[DBMS] + create Auction %v", order)
	_, err := b.db.Exec("INSERT INTO auctions (uuid, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);",
		order.UUID,
		order.PlayerID.String(),
		order.CurrencyID,
		order.Price.String(),
		order.Rnd.String(),
		order.ValidUntil.String(),
		order.Signature,
		order.State,
		order.StateExtra,
		order.Seller,
	)
	return err
}

func (b *Storage) GetOpenAuctions() ([]*Auction, error) {
	auctions, err := b.GetAuctions()
	if err != nil {
		return nil, err
	}
	var openAunction []*Auction
	for _, auction := range auctions {
		if auction.State == AUCTION_STARTED ||
			auction.State == AUCTION_ASSET_FROZEN ||
			auction.State == AUCTION_PAYING {
			openAunction = append(openAunction, auction)
		}
	}
	return openAunction, nil
}

func (b *Storage) UpdateAuctionState(uuid uuid.UUID, state AuctionState, stateExtra string) error {
	_, err := b.db.Exec("UPDATE auctions SET state=$1, state_extra=$2 WHERE uuid=$3;", state, stateExtra, uuid)
	return err
}

func (b *Storage) UpdateAuctionPaymentUrl(uuid uuid.UUID, url string) error {
	_, err := b.db.Exec("UPDATE auctions SET payment_url=$1 WHERE uuid=$2;", url, uuid)
	return err
}

func (b *Storage) GetAuctions() ([]*Auction, error) {
	var orders []*Auction
	rows, err := b.db.Query("SELECT uuid, player_id, currency_id, price, rnd, valid_until, signature, state, payment_url, state_extra, seller FROM auctions;")
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	for rows.Next() {
		var order Auction
		var playerID sql.NullString
		var price sql.NullString
		var rnd sql.NullString
		var validUntil sql.NullString
		err = rows.Scan(
			&order.UUID,
			&playerID,
			&order.CurrencyID,
			&price,
			&rnd,
			&validUntil,
			&order.Signature,
			&order.State,
			&order.PaymentURL,
			&order.StateExtra,
			&order.Seller,
		)
		if err != nil {
			return orders, err
		}
		order.PlayerID, _ = new(big.Int).SetString(playerID.String, 10)
		order.Price, _ = new(big.Int).SetString(price.String, 10)
		order.Rnd, _ = new(big.Int).SetString(rnd.String, 10)
		order.ValidUntil, _ = new(big.Int).SetString(validUntil.String, 10)
		orders = append(orders, &order)
	}
	return orders, nil
}
