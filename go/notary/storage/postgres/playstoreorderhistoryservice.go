package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *StorageHistoryTx) PlayStoreUpdateState(order storage.PlaystoreOrder) error {
	currentOrder, err := b.Tx.PlayStoreOrder(order.OrderId)
	if err != nil {
		return err
	}
	if currentOrder == nil {
		return nil
	}
	if *currentOrder == order {
		return nil
	}
	if err := b.Tx.PlayStoreUpdateState(order); err != nil {
		return err
	}
	return playStoreInsertHistory(b.Tx.tx, order)
}

func (b *StorageHistoryTx) PlayStoreInsert(order storage.PlaystoreOrder) error {
	if err := b.Tx.PlayStoreInsert(order); err != nil {
		return err
	}
	return playStoreInsertHistory(b.Tx.tx, order)
}

func playStoreInsertHistory(tx *sql.Tx, order storage.PlaystoreOrder) error {
	_, err := tx.Exec(`INSERT INTO playstore_orders_histories (
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
