package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

type StorageHistoryService struct {
	StorageService
}

type StorageHistoryTx struct {
	Tx
}

func NewStorageHistoryService(db *sql.DB) *StorageHistoryService {
	return &StorageHistoryService{*NewStorageService(db)}
}

func (b *StorageHistoryService) Begin() (storage.Tx, error) {
	var err error
	tx, err := b.db.Begin()
	if err != nil {
		return nil, err
	}
	historyTx := StorageHistoryTx{}
	historyTx.tx = tx
	return &historyTx, nil
}

func (b *StorageHistoryTx) AuctionsHistoriesCount() int {
	var count int
	b.Tx.tx.QueryRow("SELECT count(*) FROM auctions_histories;").Scan(&count)
	return count
}

func (b *StorageHistoryTx) BidsHistoriesCount() int {
	var count int
	b.Tx.tx.QueryRow("SELECT count(*) FROM bids_histories;").Scan(&count)
	return count
}

func (b *StorageHistoryTx) PlaystoreHistoriesCount() int {
	var count int
	b.Tx.tx.QueryRow("SELECT count(*) FROM playstore_orders_histories;").Scan(&count)
	return count
}

func (b *StorageHistoryTx) OffersHistoriesCount() int {
	var count int
	b.Tx.tx.QueryRow("SELECT count(*) FROM offers_histories_v2;").Scan(&count)
	return count
}
