package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

type PlaystoreOrderService struct {
	tx *sql.Tx
}

func NewPlaystoreOrderService(tx *sql.Tx) *PlaystoreOrderService {
	return &PlaystoreOrderService{
		tx: tx,
	}
}

func (b PlaystoreOrderService) UpdateState(order *storage.PlaystoreOrder) error {
	_, err := b.tx.Exec(`UPDATE playstore_orders SET 
		state=$1, 
		state_extra=$2
		WHERE order_id=$3;`,
		order.State,
		order.StateExtra,
		order.OrderId,
	)
	return err
}

func (b PlaystoreOrderService) PendingOrders() ([]storage.PlaystoreOrder, error) {
	rows, err := b.tx.Query(`SELECT 
	order_id,
	package_name,
	product_id,
	purchase_token,
	player_id,
	team_id,
	signature,
	state,
	state_extra 
	FROM playstore_orders WHERE NOT (state='failed' OR state='refunded' OR state='complete');`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []storage.PlaystoreOrder{}
	for rows.Next() {
		order := storage.PlaystoreOrder{}
		err = rows.Scan(
			&order.OrderId,
			&order.PackageName,
			&order.ProductId,
			&order.PurchaseToken,
			&order.PlayerId,
			&order.TeamId,
			&order.Signature,
			&order.State,
			&order.StateExtra,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (b PlaystoreOrderService) Order(orderId string) (*storage.PlaystoreOrder, error) {
	rows, err := b.tx.Query(`SELECT 
	package_name,
	product_id,
	purchase_token,
	player_id,
	team_id,
	signature,
	state, 
	state_extra 
	FROM playstore_orders WHERE order_id=$1;`, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	order := storage.PlaystoreOrder{}
	order.OrderId = orderId

	err = rows.Scan(
		&order.PackageName,
		&order.ProductId,
		&order.PurchaseToken,
		&order.PlayerId,
		&order.TeamId,
		&order.Signature,
		&order.State,
		&order.StateExtra,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (b PlaystoreOrderService) Insert(order *storage.PlaystoreOrder) error {
	_, err := b.tx.Exec(`INSERT INTO playstore_orders (
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
