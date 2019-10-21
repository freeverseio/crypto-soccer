package storage

import (
	"math/big"
)

type Order struct {
	Auction Auction
	Bid       Bid
}

func (b *Storage) GetOrders() ([]Order, error) {
	var orders []Order
	Auctions, err := b.GetAuctions()
	if err != nil {
		return orders, err
	}
	bids, err := b.GetBids()
	if err != nil {
		return orders, err
	}
	for _, Auction := range Auctions {
		Bid := b.findBet(bids, Auction.PlayerID)
		if Bid != nil {
			orders = append(orders, Order{
				Auction: Auction,
				Bid:       *Bid,
			})
		}
	}
	return orders, nil
}

func (b *Storage) findBet(orders []Bid, playerId *big.Int) *Bid {
	// for _, order := range orders {
	// 	if order.PlayerID.Cmp(playerId) == 0 {
	// 		return &order
	// 	}
	// }
	return nil
}
