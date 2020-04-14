package storage

import (
	"database/sql"
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
	AuctionStarted              AuctionState = "started"
	AuctionEnded                AuctionState = "ended"
	AuctionCancelled            AuctionState = "cancelled"
	AuctionFailed               AuctionState = "failed"
)

type Auction struct {
	ID         string
	PlayerID   string
	CurrencyID int
	Price      int
	Rnd        int
	ValidUntil string
	Signature  string
	State      AuctionState
	StateExtra string
	PaymentURL string
	Seller     string
}

func NewAuction() *Auction {
	auction := Auction{}
	auction.State = AuctionStarted
	return &auction
}

func AuctionByID(tx *sql.Tx, ID string) (*Auction, error) {
	rows, err := tx.Query("SELECT player_id, currency_id, price, rnd, valid_until, signature, state, payment_url, state_extra, seller FROM auctions WHERE id = $1;", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rows.Next()
	var auction Auction
	auction.ID = ID
	err = rows.Scan(
		&auction.PlayerID,
		&auction.CurrencyID,
		&auction.Price,
		&auction.Rnd,
		&auction.ValidUntil,
		&auction.Signature,
		&auction.State,
		&auction.PaymentURL,
		&auction.StateExtra,
		&auction.Seller,
	)
	return &auction, err
}

func (b Auction) Insert(tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO auctions (id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, payment_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);",
		b.ID,
		b.PlayerID,
		b.CurrencyID,
		b.Price,
		b.Rnd,
		b.ValidUntil,
		b.Signature,
		b.State,
		b.StateExtra,
		b.Seller,
		b.PaymentURL,
	)
	return err
}

// func (b *Storage) CreateAuction(order Auction) error {
// 	log.Infof("[DBMS] + create Auction %v", order)
// 	_, err := b.db.Exec("INSERT INTO auctions (id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);",
// 		order.ID,
// 		order.PlayerID.String(),
// 		order.CurrencyID,
// 		order.Price.String(),
// 		order.Rnd.String(),
// 		order.ValidUntil,
// 		order.Signature,
// 		order.State,
// 		order.StateExtra,
// 		order.Seller,
// 	)
// 	return err
// }

// func (b *Storage) GetOpenAuctions() ([]*Auction, error) {
// 	auctions, err := b.GetAuctions()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var openAunction []*Auction
// 	for _, auction := range auctions {
// 		if auction.State == AUCTION_STARTED ||
// 			auction.State == AUCTION_ASSET_FROZEN ||
// 			auction.State == AUCTION_PAYING {
// 			openAunction = append(openAunction, auction)
// 		}
// 	}
// 	return openAunction, nil
// }

// func (b *Storage) UpdateAuctionState(ID string, state AuctionState, stateExtra string) error {
// 	_, err := b.db.Exec("UPDATE auctions SET state=$1, state_extra=$2 WHERE id=$3;", state, stateExtra, ID)
// 	return err
// }

// func (b *Storage) UpdateAuctionPaymentUrl(ID string, url string) error {
// 	_, err := b.db.Exec("UPDATE auctions SET payment_url=$1 WHERE id=$2;", url, ID)
// 	return err
// }

// func GetPendingAuctions() ([]*Auction, error) {
// 	return nil, nil
// }

// func (b *Storage) GetAuctions() ([]*Auction, error) {
// 	var orders []*Auction
// 	rows, err := b.db.Query("SELECT uuid, player_id, currency_id, price, rnd, valid_until, signature, state, payment_url, state_extra, seller FROM auctions;")
// 	if err != nil {
// 		return orders, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var order Auction
// 		var playerID sql.NullString
// 		var price sql.NullString
// 		var rnd sql.NullString
// 		err = rows.Scan(
// 			&order.ID,
// 			&playerID,
// 			&order.CurrencyID,
// 			&price,
// 			&rnd,
// 			&order.ValidUntil,
// 			&order.Signature,
// 			&order.State,
// 			&order.PaymentURL,
// 			&order.StateExtra,
// 			&order.Seller,
// 		)
// 		if err != nil {
// 			return orders, err
// 		}
// 		order.PlayerID, _ = new(big.Int).SetString(playerID.String, 10)
// 		order.Price, _ = new(big.Int).SetString(price.String, 10)
// 		order.Rnd, _ = new(big.Int).SetString(rnd.String, 10)
// 		orders = append(orders, &order)
// 	}
// 	return orders, nil
// }
