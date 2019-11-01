package auctionmachine

import (
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

func OrderByDescExtraPrice(bids []storage.Bid) []storage.Bid {
	for i := 0; i < len(bids); i++ {
		for j := 0; j < len(bids)-1-i; j++ {
			if bids[j].ExtraPrice < bids[j+1].ExtraPrice {
				bids[j], bids[j+1] = bids[j+1], bids[j]
			}
		}
	}
	return bids
}
