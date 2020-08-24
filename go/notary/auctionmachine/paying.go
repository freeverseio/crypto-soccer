package auctionmachine

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/marketpay"
	"github.com/freeverseio/crypto-soccer/go/notary/bidmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *AuctionMachine) ProcessPaying(market marketpay.MarketPayService) error {
	if err := b.checkState(storage.AuctionPaying); err != nil {
		return err
	}

	bid := bidmachine.FirstAlive(b.bids)
	if bid == nil {
		b.SetState(storage.AuctionFailed, "No available healty bid")
		return nil
	}

	bidMachine, err := bidmachine.New(
		market,
		b.auction,
		bid,
		b.contracts,
		b.freeverse,
	)
	if err != nil {
		return err
	}

	err = bidMachine.Process()
	if err != nil {
		return err
	}
	if bid.State == storage.BidPaid {
		if err := b.transferAuction(*bid); err != nil {
			return err
		}
		order, err := market.GetOrder(bid.PaymentID)
		if err != nil {
			return err
		}
		b.auction.PaymentURL = order.SettlorShortlink.ShortURL
		bid.State = storage.BidPaid

		b.SetState(storage.AuctionWithdrableBySeller, "")
	}

	return nil
}

func (b AuctionMachine) transferAuction(bid storage.Bid) error {
	// transfer the auction
	bidHiddenPrice, err := signer.BidHiddenPrice(b.contracts.Market, big.NewInt(bid.ExtraPrice), big.NewInt(bid.Rnd))
	if err != nil {
		return err
	}
	auctionHiddenPrice, err := signer.HashPrivateMsg(uint8(b.auction.CurrencyID), big.NewInt(b.auction.Price), big.NewInt(b.auction.Rnd))
	if err != nil {
		return err
	}
	playerId, _ := new(big.Int).SetString(b.auction.PlayerID, 10)
	if playerId == nil {
		return errors.New("invalid playerid")
	}
	teamId, _ := new(big.Int).SetString(bid.TeamID, 10)
	if playerId == nil {
		return errors.New("invalid teamid")
	}

	isOffer := b.offer != nil

	var validUntil int64
	if isOffer {
		validUntil = b.offer.ValidUntil
	} else {
		validUntil = b.auction.ValidUntil
	}

	var sig [2][32]byte
	var sigV uint8
	_, err = signer.HashBidMessage(
		b.contracts.Market,
		uint8(b.auction.CurrencyID),
		big.NewInt(b.auction.Price),
		big.NewInt(b.auction.Rnd),
		b.auction.ValidUntil,
		playerId,
		big.NewInt(bid.ExtraPrice),
		big.NewInt(bid.Rnd),
		teamId,
		isOffer,
	)
	if err != nil {
		return err
	}
	sig[0], sig[1], sigV, err = signer.RSV(bid.Signature)
	if err != nil {
		return err
	}
	auth := bind.NewKeyedTransactor(b.freeverse)
	auth.GasPrice = big.NewInt(1000000000) // in xdai is fixe to 1 GWei
	tx, err := b.contracts.Market.CompletePlayerAuction(
		auth,
		auctionHiddenPrice,
		big.NewInt(validUntil),
		playerId,
		bidHiddenPrice,
		teamId,
		sig,
		sigV,
		isOffer,
	)
	if err != nil {
		b.SetState(storage.AuctionWithdrableByBuyer, err.Error())
		return err
	}
	receipt, err := helper.WaitReceipt(b.contracts.Client, tx, 60)
	if err != nil {
		b.SetState(storage.AuctionWithdrableByBuyer, "Timeout waiting for the receipt")
		return err
	}
	if receipt.Status == 0 {
		b.SetState(storage.AuctionWithdrableByBuyer, "Mined but receipt.Status == 0")
		return err
	}
	return nil
}
