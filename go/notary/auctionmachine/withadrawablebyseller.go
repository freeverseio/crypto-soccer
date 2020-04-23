package auctionmachine

import (
	"fmt"

	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *AuctionMachine) ProcessWithdrawableBySeller(market marketpay.IMarketPay) error {
	if b.State() != storage.AuctionWithdrableBySeller {
		return fmt.Errorf("Wrong state %v", b.State())
	}

	return nil
}
