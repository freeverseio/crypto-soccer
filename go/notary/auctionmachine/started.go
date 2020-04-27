package auctionmachine

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *AuctionMachine) processStarted() error {
	if b.auction.State != storage.AuctionStarted {
		return errors.New("Started: wrong state")
	}

	// check if expired
	now := time.Now().Unix()
	if now > b.auction.ValidUntil {
		b.SetState(storage.AuctionEnded, "expired")
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
		b.SetState(storage.AuctionFailed, "auction is frozen in other market")
		return nil
	}

	// check if seller is the owner
	owner, err := b.contracts.Market.GetOwnerPlayer(&bind.CallOpts{}, playerId)
	if err != nil {
		return err
	}
	if owner.String() != b.auction.Seller {
		b.SetState(storage.AuctionFailed, fmt.Sprintf("seller %s is not the owner %s", b.auction.Seller, owner.String()))
		return nil
	}

	// if no bids do nothing
	if len(b.bids) == 0 {
		return nil
	}

	// if has bids let's freeze it
	auctionHiddenPrice, err := signer.HashPrivateMsg(
		uint8(b.auction.CurrencyID),
		big.NewInt(b.auction.Price),
		big.NewInt(b.auction.Rnd),
	)
	if err != nil {
		return err
	}
	var sig [2][32]byte
	var sigV uint8
	sig[0], sig[1], sigV, err = signer.RSV(b.auction.Signature)
	if err != nil {
		return err
	}
	auth := bind.NewKeyedTransactor(b.freeverse)
	auth.GasPrice = big.NewInt(1000000000) // in xdai is fixe to 1 GWei
	tx, err := b.contracts.Market.FreezePlayer(
		auth,
		auctionHiddenPrice,
		big.NewInt(b.auction.ValidUntil),
		playerId,
		sig,
		sigV,
	)
	if err != nil {
		b.SetState(storage.AuctionFailed, "Failed to freeze: "+err.Error())
		return err
	}
	receipt, err := helper.WaitReceipt(b.contracts.Client, tx, 60)
	if err != nil {
		b.SetState(storage.AuctionFailed, "Failed to Freeze: waiting for receipt timeout")
		return err
	}
	if receipt.Status == 0 {
		b.SetState(storage.AuctionFailed, "Failed to Freeze: mined but receipt status is failed")
		return err
	}

	b.auction.State = storage.AuctionAssetFrozen
	return nil
}
