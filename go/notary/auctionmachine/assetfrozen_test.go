package auctionmachine_test

import (
	"testing"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestAssetFrozen(t *testing.T) {
	auction := storage.NewAuction()
	auction.ID = "f1d4501c5158a9018b1618ec4d471c66b663d8f6bffb6e70a0c6584f5c1ea94a"
	auction.ValidUntil = time.Now().Unix() + 100
	auction.PlayerID = "274877906944"
	auction.CurrencyID = 1
	auction.Price = 41234
	auction.Rnd = 4232
	auction.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e"
	auction.Signature = "381bf58829e11790830eab9924b123d1dbe96dd37b10112729d9d32d476c8d5762598042bb5d5fd63f668455aa3a2ce4e2632241865c26ababa231ad212b5f151b"

	auction.State = storage.AuctionStarted
	shouldQueryMarketPay := true
	m, err := auctionmachine.New(*auction, nil, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.Error(t, m.ProcessAssetFrozen(), "auction[f1d4501c5158a9018b1618ec4d471c66b663d8f6bffb6e70a0c6584f5c1ea94a|started] is not in state asset_frozen")

	auction.State = storage.AuctionAssetFrozen
	m, err = auctionmachine.New(*auction, nil, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessAssetFrozen())
}
