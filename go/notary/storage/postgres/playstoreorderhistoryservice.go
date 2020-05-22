package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

type PlaystoreOrderHistoryService struct {
	tx      *sql.Tx
	service storage.PlaystoreOrderService
}

func NewPlaystoreOrderHistoryService(tx *sql.Tx, service storage.PlaystoreOrderService) *PlaystoreOrderHistoryService {
	return &PlaystoreOrderHistoryService{
		tx:      tx,
		service: service,
	}
}

func (b PlaystoreOrderHistoryService) UpdateState(order storage.PlaystoreOrder) error {
	if err := b.service.UpdateState(order); err != nil {
		return err
	}
	return b.insertHistory(order)
}

func (b PlaystoreOrderHistoryService) PendingOrders() ([]storage.PlaystoreOrder, error) {
	return b.service.PendingOrders()
}

func (b PlaystoreOrderHistoryService) Order(orderId string) (*storage.PlaystoreOrder, error) {
	return b.service.Order(orderId)
}

func (b PlaystoreOrderHistoryService) Insert(order storage.PlaystoreOrder) error {
	if err := b.service.Insert(order); err != nil {
		return err
	}
	return b.insertHistory(order)
}

func (b PlaystoreOrderHistoryService) insertHistory(order storage.PlaystoreOrder) error {
	_, err := b.tx.Exec(`INSERT INTO playstore_orders_histories (
		order_id, 
		package_name,
		product_id,
		purchase_token,
		player_id,
		team_id,
		signature,
		state, 
		state_extra
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		order.OrderId,
		order.PackageName,
		order.ProductId,
		order.PurchaseToken,
		order.PlayerId,
		order.TeamId,
		order.Signature,
		order.State,
		order.StateExtra,
	)
	return err
}
