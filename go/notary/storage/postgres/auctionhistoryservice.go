package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

type AuctionHistoryService struct {
	AuctionService
}

func NewAuctionHistoryService(tx *sql.Tx) *AuctionHistoryService {
	return &AuctionHistoryService{*NewAuctionService(tx)}
}

func (b AuctionHistoryService) Bid() storage.BidService {
	return *NewBidHistoryService(b.tx)
}
