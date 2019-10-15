package storage

import (
	"math/big"
)

type Order struct {
	SellOrder SellOrder
	Bet       Bet
}

func (b *Storage) GetOrders() ([]Order, error) {
	var orders []Order
	sellOrders, err := b.GetSellOrders()
	if err != nil {
		return orders, err
	}
	Bets, err := b.GetBets()
	if err != nil {
		return orders, err
	}
	for _, sellOrder := range sellOrders {
		Bet := b.findBet(Bets, sellOrder.PlayerID)
		if Bet != nil {
			orders = append(orders, Order{
				SellOrder: sellOrder,
				Bet:       *Bet,
			})
		}
	}
	return orders, nil
}

func (b *Storage) findBet(orders []Bet, playerId *big.Int) *Bet {
	// for _, order := range orders {
	// 	if order.PlayerID.Cmp(playerId) == 0 {
	// 		return &order
	// 	}
	// }
	return nil
}
