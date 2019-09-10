package storage

import (
	"math/big"

	log "github.com/sirupsen/logrus"
)

type Order struct {
	SellOrder SellOrder
	BuyOrder  BuyOrder
}

func (b *Storage) GetOrders() ([]Order, error) {
	var orders []Order
	sellOrders, err := b.GetSellOrders()
	if err != nil {
		return orders, err
	}
	buyOrders, err := b.GetBuyOrders()
	if err != nil {
		return orders, err
	}
	for _, sellOrder := range sellOrders {
		buyOrder := b.findBuyOrder(buyOrders, sellOrder.PlayerId)
		if buyOrder != nil {
			orders = append(orders, Order{
				SellOrder: sellOrder,
				BuyOrder:  *buyOrder,
			})
		}
	}
	return orders, nil
}

func (b *Storage) findBuyOrder(orders []BuyOrder, playerId *big.Int) *BuyOrder {
	for _, order := range orders {
		if order.PlayerId == playerId.Uint64() {
			return &order
		}
	}
	return nil
}

func (b *Storage) DeleteOrder(playerId *big.Int) error {
	log.Infof("[DBMS] - delete order %v", playerId)
	err := b.DeleteBuyOrder(playerId.Uint64())
	if err != nil {
		return err
	}
	err = b.DeleteSellOrder(playerId)
	return err
}
