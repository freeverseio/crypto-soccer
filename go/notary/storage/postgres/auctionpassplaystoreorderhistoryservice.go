package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *StorageHistoryTx) AuctionPassPlayStoreUpdateState(order storage.AuctionPassPlaystoreOrder) error {
	currentOrder, err := b.Tx.AuctionPassPlayStoreOrder(order.OrderId)
	if err != nil {
		return err
	}
	if currentOrder == nil {
		return nil
	}
	if *currentOrder == order {
		return nil
	}
	if err := b.Tx.AuctionPassPlayStoreUpdateState(order); err != nil {
		return err
	}
	return auctionPassPlayStoreInsertHistory(b.Tx.tx, order)
}

func (b *StorageHistoryTx) AuctionPassPlayStoreInsert(order storage.AuctionPassPlaystoreOrder) error {
	if err := b.Tx.AuctionPassPlayStoreInsert(order); err != nil {
		return err
	}
	return auctionPassPlayStoreInsertHistory(b.Tx.tx, order)
}

func auctionPassPlayStoreInsertHistory(tx *sql.Tx, order storage.AuctionPassPlaystoreOrder) error {
	_, err := tx.Exec(`INSERT INTO playstore_orders_histories (
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
