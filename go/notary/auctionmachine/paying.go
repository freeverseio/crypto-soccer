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
	playerId, _ := new(big.Int).SetString(b.auction.PlayerID, 10)
	if playerId == nil {
		return errors.New("invalid playerid")
	}
	auctionId, err := signer.ComputeAuctionId(
		uint8(b.auction.CurrencyID),
		big.NewInt(b.auction.Price),
		big.NewInt(b.auction.Rnd),
		b.auction.ValidUntil,
		b.auction.OfferValidUntil,
		playerId,
	)
	if err != nil {
		return errors.Wrap(err, "AuctionHiddenPrice")
	}
	teamId, _ := new(big.Int).SetString(bid.TeamID, 10)
	if playerId == nil {
		return errors.New("invalid teamid")
	}

	var sig [2][32]byte
	var sigV uint8
	_, err = signer.HashBidMessageFromAuctionId(
		b.contracts.Market,
		auctionId,
		big.NewInt(bid.ExtraPrice),
		big.NewInt(bid.Rnd),
		teamId,
	)
	if err != nil {
		return errors.Wrap(err, "HashBidMessageFromAuctionId")
	}
	sig[0], sig[1], sigV, err = signer.RSV(bid.Signature)
	if err != nil {
		return errors.Wrap(err, "RSV")
	}
	auth := bind.NewKeyedTransactor(b.freeverse)
	auth.GasPrice = big.NewInt(1000000000) // in xdai is fixe to 1 GWei
	// fails here
	tx, err := b.contracts.Market.CompletePlayerAuction(
		auth,
		auctionId,
		playerId,
		bidHiddenPrice,
		teamId,
		sig,
		sigV,
	)
	if err != nil {
		return errors.Wrapf(err, "CompletePlayerAuction auctionId: %v, validUntil %v, offerValidUntil %v, playerId, %v, bidHiddenPrice: 0x%v, teamId: %v, sig[0]: 0x%v, sig[1]: 0x%v , sigV: %v",
			auctionId.Hex(),
			big.NewInt(b.auction.ValidUntil),
			big.NewInt(b.auction.OfferValidUntil),
			playerId,
			hex.EncodeToString(bidHiddenPrice[:]),
			teamId,
			hex.EncodeToString(sig[0][:]),
			hex.EncodeToString(sig[1][:]),
			sigV,
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
