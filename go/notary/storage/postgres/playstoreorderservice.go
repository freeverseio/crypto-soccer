package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b StorageService) PlayStoreUpdateState(tx *sql.Tx, order storage.PlaystoreOrder) error {
	_, err := tx.Exec(`UPDATE playstore_orders SET 
		state=$1, 
		state_extra=$2
		WHERE order_id=$3;`,
		order.State,
		order.StateExtra,
		order.OrderId,
	)
	return err
}

func (b StorageService) PlayStorePendingOrders(tx *sql.Tx) ([]storage.PlaystoreOrder, error) {
	rows, err := tx.Query(`SELECT 
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

func (b StorageService) PlayStorePendingOrdersByPlayerId(tx *sql.Tx, playerId string) ([]storage.PlaystoreOrder, error) {
	rows, err := tx.Query(`SELECT 
	order_id,
	package_name,
	product_id,
	purchase_token,
	player_id,
	team_id,
	signature,
	state,
	state_extra 
	FROM playstore_orders WHERE (player_id=$1) AND NOT (state='failed' OR state='refunded' OR state='complete');`, playerId)
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

func (b StorageService) PlayStoreOrder(tx *sql.Tx, orderId string) (*storage.PlaystoreOrder, error) {
	rows, err := tx.Query(`SELECT 
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

func (b StorageService) PlayStoreInsert(tx *sql.Tx, order storage.PlaystoreOrder) error {
	_, err := tx.Exec(`INSERT INTO playstore_orders (
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
