package auctionmachine

import (
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
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

func HasAccepted(bids []storage.Bid) bool {
	return IndexOfFirstAccepted(bids) != -1
}

func HasPaying(bids []storage.Bid) bool {
	return IndexOfFirstPaying(bids) != -1
}

func IndexOfFirstPaying(bids []storage.Bid) int {
	for idx, bid := range bids {
		if bid.State == storage. BIDPAYING {
			return idx
		}
	}
	return -1
}

func IndexOfFirstAccepted(bids []storage.Bid) int {
	for idx, bid := range bids {
		if bid.State == storage.BIDACCEPTED {
			return idx
		}
	}
	return -1
}
