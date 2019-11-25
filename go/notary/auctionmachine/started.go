package auctionmachine

import (
	"errors"
	"time"

	"github.com/freeverseio/crypto-soccer/go/helper"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func (m *AuctionMachine) processStarted() error {
	if m.Auction.State != storage.AUCTION_STARTED {
		return errors.New("Started: wrong state")
	}
	now := time.Now().Unix()

	if len(m.Bids) == 0 {
		if now > m.Auction.ValidUntil.Int64() {
			log.Infof("Auction %v STARTED -> %v", m.Auction.UUID, m.Auction.State)
			m.Auction.State = storage.AUCTION_NO_BIDS
		}
		return nil
	}

	// TODO trying to freeze the asset
	auctionHiddenPrice, err := m.signer.HashPrivateMsg(
		m.Auction.CurrencyID,
		m.Auction.Price,
		m.Auction.Rnd,
	)
	if err != nil {
		return err
	}
	var sig [3][32]byte
	var sigV uint8
	sig[0], err = m.signer.HashSellMessage(
		m.Auction.CurrencyID,
		m.Auction.Price,
		m.Auction.Rnd,
		m.Auction.ValidUntil,
		m.Auction.PlayerID,
	)
	if err != nil {
		return err
	}
	sig[1], sig[2], sigV, err = m.signer.RSV(m.Auction.Signature)
	if err != nil {
		return err
	}
	tx, err := m.contracts.Market.FreezePlayer(
		bind.NewKeyedTransactor(m.freeverse),
		auctionHiddenPrice,
		m.Auction.ValidUntil,
		m.Auction.PlayerID,
		sig,
		sigV,
	)
	if err != nil {
		log.Error(err)
		m.Auction.State = storage.AUCTION_FAILED
		m.Auction.StateExtra = "Failed to freeze: " + err.Error()
		return nil
	}
	receipt, err := helper.WaitReceipt(m.contracts.Client, tx, 60)
	if err != nil {
		log.Error("Timeout waiting receipt for freeze")
		m.Auction.State = storage.AUCTION_FAILED
		m.Auction.State = "Failed to Freeze: waiting for receipt timeout"
		return nil
	}
	if receipt.Status == 0 {
		log.Error("Freeze mined but failed")
		m.Auction.State = storage.AUCTION_FAILED
		m.Auction.State = "Failed to Freeze: mined but receipt status is failed"
		return nil
	}

	log.Infof("[auction] %v STARTER -> ASSET_FROZEN", m.Auction.UUID)
	m.Auction.State = storage.AUCTION_ASSET_FROZEN
	return nil
}
