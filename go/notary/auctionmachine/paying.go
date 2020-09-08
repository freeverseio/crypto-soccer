package auctionmachine

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/marketpay"
	"github.com/freeverseio/crypto-soccer/go/notary/bidmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
			b.SetState(storage.AuctionWithdrableByBuyer, err.Error())
			return nil
		}
		order, err := market.GetOrder(bid.PaymentID)
		if err != nil {
			return err
		}
		b.auction.PaymentURL = order.SettlorShortlink.ShortURL

		b.SetState(storage.AuctionWithdrableBySeller, "")
	}

	return nil
}

func (b AuctionMachine) transferAuction(bid storage.Bid) error {
	log.Debugf("Transfer player %v to team %v", b.auction.PlayerID, bid.TeamID)

	// transfer the auction
	bidHiddenPrice, err := signer.BidHiddenPrice(b.contracts.Market, big.NewInt(bid.ExtraPrice), big.NewInt(bid.Rnd))
	if err != nil {
		return errors.Wrap(err, "BidHiddenPrice")
	}
	auctionHiddenPrice, err := signer.HashPrivateMsg(uint8(b.auction.CurrencyID), big.NewInt(b.auction.Price), big.NewInt(b.auction.Rnd))
	if err != nil {
		return errors.Wrap(err, "AuctionHiddenPrice")
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
		return errors.Wrap(err, "HashBidMessage")
	}
	sig[0], sig[1], sigV, err = signer.RSV(bid.Signature)
	if err != nil {
		return errors.Wrap(err, "RSV")
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
		return errors.Wrapf(err, "CompletePlayerAuction auctionHiddenPrice: %v, validUntil %v, playerId, %v, bidHiddenPrice: 0x%v, teamId: %v, sig[0]: 0x%v, sig[1]: 0x%v , sigV: %v, isOffer: %v",
			auctionHiddenPrice.Hex(),
			big.NewInt(validUntil),
			playerId,
			hex.EncodeToString(bidHiddenPrice[:]),
			teamId,
			hex.EncodeToString(sig[0][:]),
			hex.EncodeToString(sig[1][:]),
			sigV,
			isOffer,
		)
	}
	receipt, err := helper.WaitReceipt(b.contracts.Client, tx, 60)
	if err != nil {
		return errors.Wrap(err, "WaitReceipt")
	}
	if receipt.Status == 0 {
		return errors.New("Status != 0")
	}
	return nil
}
