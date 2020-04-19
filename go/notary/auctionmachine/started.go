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
	if b.auction.State != storage.AuctionStarted {
		return errors.New("Started: wrong state")
	}

	// check if expired
	now := time.Now().Unix()
	if now > b.auction.ValidUntil {
		b.auction.State = storage.AuctionEnded
		return nil
	}

	playerId, _ := new(big.Int).SetString(b.auction.PlayerID, 10)
	if playerId == nil {
		return fmt.Errorf("Invalid PlayerId %x", b.auction.PlayerID)
	}

	// check if is frozen
	isFrozen, err := b.contracts.Market.IsPlayerFrozenInAnyMarket(&bind.CallOpts{}, playerId)
	if err != nil {
		return err
	}
	if isFrozen {
		b.auction.State = storage.AuctionFailed
		b.auction.StateExtra = "auction is frozen in other market"
		return nil
	}

	// check if seller is the owner
	owner, err := b.contracts.Market.GetOwnerPlayer(&bind.CallOpts{}, playerId)
	if err != nil {
		return err
	}
	if owner.String() != b.auction.Seller {
		b.auction.State = storage.AuctionFailed
		b.auction.StateExtra = fmt.Sprintf("seller %s is not the owner %s", b.auction.Seller, owner.String())
		return nil
	}

	// 	// TODO trying to freeze the asset
	// 	auctionHiddenPrice, err := signer.HashPrivateMsg(
	// 		m.auction.CurrencyID,
	// 		m.auction.Price,
	// 		m.auction.Rnd,
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	var sig [2][32]byte
	// 	var sigV uint8
	// 	_, err = signer.HashSellMessage(
	// 		m.auction.CurrencyID,
	// 		m.auction.Price,
	// 		m.auction.Rnd,
	// 		m.auction.ValidUntil,
	// 		m.auction.PlayerID,
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	sig[0], sig[1], sigV, err = signer.RSV(m.auction.Signature)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	tx, err := m.contracts.Market.FreezePlayer(
	// 		bind.NewKeyedTransactor(m.freeverse),
	// 		auctionHiddenPrice,
	// 		big.NewInt(m.auction.ValidUntil),
	// 		m.auction.PlayerID,
	// 		sig,
	// 		sigV,
	// 	)
	// 	if err != nil {
	// 		log.Error(err)
	// 		m.auction.State = storage.AUCTION_FAILED
	// 		m.auction.StateExtra = "Failed to freeze: " + err.Error()
	// 		return nil
	// 	}
	// 	receipt, err := helper.WaitReceipt(m.contracts.Client, tx, 60)
	// 	if err != nil {
	// 		log.Error("Timeout waiting receipt for freeze")
	// 		m.auction.State = storage.AUCTION_FAILED
	// 		m.auction.State = "Failed to Freeze: waiting for receipt timeout"
	// 		return nil
	// 	}
	// 	if receipt.Status == 0 {
	// 		log.Error("Freeze mined but failed")
	// 		m.auction.State = storage.AUCTION_FAILED
	// 		m.auction.State = "Failed to Freeze: mined but receipt status is failed"
	// 		return nil
	// 	}

	// 	log.Infof("[auction] %v STARTER -> ASSET_FROZEN", m.auction.UUID)
	// 	m.auction.State = storage.AUCTION_ASSET_FROZEN
	return nil
}
