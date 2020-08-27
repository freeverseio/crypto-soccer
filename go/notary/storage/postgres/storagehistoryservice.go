package postgres

import (
	"database/sql"
)

type StorageHistoryService struct {
	StorageService
}

func NewStorageHistoryService(db *sql.DB) *StorageHistoryService {
	return &StorageHistoryService{*NewStorageService(db)}
}

func (b *StorageHistoryService) AuctionsHistoriesCount() int {
	var count int
	b.StorageService.tx.QueryRow("SELECT count(*) FROM auctions_histories;").Scan(&count)
	return count
}

func (b *StorageHistoryService) BidsHistoriesCount() int {
	var count int
	b.StorageService.tx.QueryRow("SELECT count(*) FROM bids_histories;").Scan(&count)
	return count
}

func (b *StorageHistoryService) PlaystoreHistoriesCount() int {
	var count int
	b.StorageService.tx.QueryRow("SELECT count(*) FROM playstore_orders_histories;").Scan(&count)
	return count
}

func (b *StorageHistoryService) OffersHistoriesCount() int {
	var count int
	b.StorageService.tx.QueryRow("SELECT count(*) FROM offers_histories;").Scan(&count)
	return count
}
