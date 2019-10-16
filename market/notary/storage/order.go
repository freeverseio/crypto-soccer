package storage

import (
	"math/big"
)

type Order struct {
	SellOrder SellOrder
	Bid       Bid
}

func (b *Storage) GetOrders() ([]Order, error) {
	var orders []Order
	sellOrders, err := b.GetSellOrders()
	if err != nil {
		return orders, err
	}
	bids, err := b.Getbids()
	if err != nil {
		return orders, err
	}
	for _, sellOrder := range sellOrders {
		Bid := b.findBet(bids, sellOrder.PlayerID)
		if Bid != nil {
			orders = append(orders, Order{
				SellOrder: sellOrder,
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
