package auctionmachine

import (
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
	log "github.com/sirupsen/logrus"
)

type Paying struct {
}

func NewPaying() State {
	return &Paying{}
}

func (b *Paying) Process(m *AuctionMachine) error {
	if m.Auction.State != storage.AUCTION_PAYING {
		return errors.New("Paying: wrong state")
	}

	now := time.Now().Unix()
	if (now - m.Auction.ValidUntil.Int64()) > 60 {
		bid := m.Bids[0]
		isOffer2StartAuction := false
		bidHiddenPrice, err := m.signer.BidHiddenPrice(big.NewInt(bid.ExtraPrice), big.NewInt(bid.Rnd))
		if err != nil {
			return err
		}
		auctionHiddenPrice, err := m.signer.HashPrivateMsg(m.Auction.CurrencyID, m.Auction.Price, m.Auction.Rnd)
		if err != nil {
			return err
		}
		var sig [3][32]byte
		var sigV uint8
		sig[0], err = m.signer.HashBidMessage(
			m.Auction.CurrencyID,
			m.Auction.Price,
			m.Auction.Rnd,
			m.Auction.ValidUntil,
			m.Auction.PlayerID,
			big.NewInt(bid.ExtraPrice),
			big.NewInt(bid.Rnd),
			bid.TeamID,
			isOffer2StartAuction,
		)
		if err != nil {
			return err
		}
		sig[1], sig[2], sigV, err = m.signer.RSV(bid.Signature)
		if err != nil {
			return err
		}
		tx, err := m.market.CompletePlayerAuction(
			bind.NewKeyedTransactor(m.freeverse),
			bidHiddenPrice,
			m.Auction.ValidUntil,
			m.Auction.PlayerID,
			auctionHiddenPrice,
			bid.TeamID,
			sig,
			sigV,
			isOffer2StartAuction,
		)
		if err != nil {
			log.Error(err)
			m.Auction.State = storage.AUCTION_FAILED_TO_PAY
			m.SetState(NewFailedToPay())
			return nil
		}
		receipt, err := helper.WaitReceipt(m.client, tx, 60)
		if err != nil {
			log.Error(err)
			m.Auction.State = storage.AUCTION_FAILED_TO_PAY
			m.SetState(NewFailedToPay())
			return nil
		}
		if receipt.Status == 0 {
			log.Error("Complete mined but failed")
			m.Auction.State = storage.AUCTION_FAILED_TO_PAY
			m.SetState(NewFailedToPay())
			return nil
		}

		m.Auction.State = storage.AUCTION_PAID
		m.SetState(NewPaid())
	}

	return nil
}
