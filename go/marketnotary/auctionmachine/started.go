package auctionmachine

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type Started struct {
}

func NewStarted() State {
	return &Started{}
}

func (b *Started) Process(m *AuctionMachine) error {
	if m.Auction.State != storage.AUCTION_STARTED {
		return errors.New("Started: wrong state")
	}
	now := time.Now().Unix()

	if len(m.Bids) == 0 {
		if now > m.Auction.ValidUntil.Int64() {
			m.Auction.State = storage.AUCTION_NO_BIDS
			m.SetState(NewNoBids())
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
	_, err = m.market.FreezePlayer(
		bind.NewKeyedTransactor(m.freeverse),
		auctionHiddenPrice,
		m.Auction.ValidUntil,
		m.Auction.PlayerID,
		sig,
		sigV,
	)
	if err != nil {
		m.Auction.State = storage.AUCTION_FAILED_TO_FREEZE
		m.SetState(NewFailedToFreeze())
		return nil
	}

	m.Auction.State = storage.AUCTION_ASSET_FROZEN
	m.SetState(NewAssetFrozen())
	return nil
}
