package storage

import (
	"database/sql"
	"math/big"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
}

func (b *Storage) CreateAuction(order Auction) error {
	log.Infof("[DBMS] + create sell order %v", order)
	_, err := b.db.Exec("INSERT INTO auctions (uuid, player_id, currency_id, price, rnd, valid_until, signature, state) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);",
		order.UUID,
		order.PlayerID.String(),
		order.CurrencyID,
		order.Price.String(),
		order.Rnd.String(),
		order.ValidUntil.String(),
		order.Signature,
		order.State,
	)
	return err
}

func (b *Storage) GetOpenAuctions() ([]Auction, error) {
	auctions, err := b.GetAuctions()
	if err != nil {
		return nil, err
	}
	var openAunction []Auction
	for _, auction := range auctions {
		if auction.State == "STARTED" ||
			auction.State == "ASSET_FROZEN" ||
			auction.State == "PAYING" {
			openAunction = append(openAunction, auction)
		}
	}
	return openAunction, nil
}

func (b *Storage) UpdateAuctionState(auction Auction) error {
	_, err := b.db.Exec("UPDATE auctions SET state=$1 WHERE uuid=$2;", auction.State, auction.UUID)
	return err
}

func (b *Storage) GetAuctions() ([]Auction, error) {
	var orders []Auction
	rows, err := b.db.Query("SELECT uuid, player_id, currency_id, price, rnd, valid_until, signature, state FROM auctions;")
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
		)
		if err != nil {
			return orders, err
		}
		order.PlayerID, _ = new(big.Int).SetString(playerID.String, 10)
		order.Price, _ = new(big.Int).SetString(price.String, 10)
		order.Rnd, _ = new(big.Int).SetString(rnd.String, 10)
		order.ValidUntil, _ = new(big.Int).SetString(validUntil.String, 10)
		orders = append(orders, order)
	}
	return orders, nil
}
