package postgres

import (
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *Tx) AuctionPassPlayStoreUpdateState(order storage.AuctionPassPlaystoreOrder) error {
	_, err := b.tx.Exec(`UPDATE auction_pass_playstore_orders SET 
		state=$1, 
		state_extra=$2
		WHERE order_id=$3;`,
		order.State,
		order.StateExtra,
		order.OrderId,
	)
	return err
}

func (b *Tx) AuctionPassPlayStorePendingOrders() ([]storage.AuctionPassPlaystoreOrder, error) {
	rows, err := b.tx.Query(`SELECT 
	order_id,
	package_name,
	product_id,
	purchase_token,
	team_id,
	owner,
	signature,
	state,
	state_extra 
	FROM auction_pass_playstore_orders WHERE NOT (state='failed' OR state='refunded' OR state='complete');`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []storage.AuctionPassPlaystoreOrder{}
	for rows.Next() {
		order := storage.AuctionPassPlaystoreOrder{}
		err = rows.Scan(
			&order.OrderId,
			&order.PackageName,
			&order.ProductId,
			&order.PurchaseToken,
			&order.TeamId,
			&order.Owner,
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

func (b *Tx) AuctionPassPlayStoreOrder(orderId string) (*storage.AuctionPassPlaystoreOrder, error) {
	rows, err := b.tx.Query(`SELECT 
	package_name,
	product_id,
	purchase_token,
	team_id,
	owner,
	signature,
	state, 
	state_extra 
	FROM auction_pass_playstore_orders WHERE order_id=$1;`, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	order := storage.AuctionPassPlaystoreOrder{}
	order.OrderId = orderId

	err = rows.Scan(
		&order.PackageName,
		&order.ProductId,
		&order.PurchaseToken,
		&order.TeamId,
		&order.Owner,
		&order.Signature,
		&order.State,
		&order.StateExtra,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (b *Tx) AuctionPassPlayStoreInsert(order storage.AuctionPassPlaystoreOrder) error {
	_, err := b.tx.Exec(`INSERT INTO auction_pass_playstore_orders (
		order_id, 
		package_name,
		product_id,
		purchase_token,
		team_id,
		owner,
		signature,
		state, 
		state_extra
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		order.OrderId,
		order.PackageName,
		order.ProductId,
		order.PurchaseToken,
		order.TeamId,
		order.Owner,
		order.Signature,
		order.State,
		order.StateExtra,
	)
	return err
}
