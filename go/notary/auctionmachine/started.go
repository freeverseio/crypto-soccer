package auctionmachine

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *AuctionMachine) processStarted() error {
	if b.Auction.State != storage.AuctionStarted {
		return errors.New("Started: wrong state")
	}
	now := time.Now().Unix()

	if now > b.Auction.ValidUntil {
		b.Auction.State = storage.AuctionEnded
		return nil
	}

	playerId, _ := new(big.Int).SetString(b.Auction.PlayerID, 10)
	if playerId == nil {
		return fmt.Errorf("Invalid PlayerId %x", b.Auction.PlayerID)
	}
	owner, err := b.contracts.Market.GetOwnerPlayer(&bind.CallOpts{}, playerId)
	if err != nil {
		return err
	}
	if owner.String() != b.Auction.Seller {
		b.Auction.State = storage.AuctionFailed
		b.Auction.StateExtra = fmt.Sprintf("seller %s is not the owner %s", b.Auction.Seller, owner.String())
		return nil
	}

	// 	// TODO trying to freeze the asset
	// 	auctionHiddenPrice, err := signer.HashPrivateMsg(
	// 		m.Auction.CurrencyID,
	// 		m.Auction.Price,
	// 		m.Auction.Rnd,
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	var sig [2][32]byte
	// 	var sigV uint8
	// 	_, err = signer.HashSellMessage(
	// 		m.Auction.CurrencyID,
	// 		m.Auction.Price,
	// 		m.Auction.Rnd,
	// 		m.Auction.ValidUntil,
	// 		m.Auction.PlayerID,
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	sig[0], sig[1], sigV, err = signer.RSV(m.Auction.Signature)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	tx, err := m.contracts.Market.FreezePlayer(
	// 		bind.NewKeyedTransactor(m.freeverse),
	// 		auctionHiddenPrice,
	// 		big.NewInt(m.Auction.ValidUntil),
	// 		m.Auction.PlayerID,
	// 		sig,
	// 		sigV,
	// 	)
	// 	if err != nil {
	// 		log.Error(err)
	// 		m.Auction.State = storage.AUCTION_FAILED
	// 		m.Auction.StateExtra = "Failed to freeze: " + err.Error()
	// 		return nil
	// 	}
	// 	receipt, err := helper.WaitReceipt(m.contracts.Client, tx, 60)
	// 	if err != nil {
	// 		log.Error("Timeout waiting receipt for freeze")
	// 		m.Auction.State = storage.AUCTION_FAILED
	// 		m.Auction.State = "Failed to Freeze: waiting for receipt timeout"
	// 		return nil
	// 	}
	// 	if receipt.Status == 0 {
	// 		log.Error("Freeze mined but failed")
	// 		m.Auction.State = storage.AUCTION_FAILED
	// 		m.Auction.State = "Failed to Freeze: mined but receipt status is failed"
	// 		return nil
	// 	}

	// 	log.Infof("[auction] %v STARTER -> ASSET_FROZEN", m.Auction.UUID)
	// 	m.Auction.State = storage.AUCTION_ASSET_FROZEN
	return nil
}
