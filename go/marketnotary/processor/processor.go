package processor

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"time"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/auctionmachine"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go/contracts/market"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

type Processor struct {
	db        *storage.Storage
	client    *ethclient.Client
	assets    *market.Market
	freeverse *ecdsa.PrivateKey
	signer    *Signer
}

func NewProcessor(db *storage.Storage, ethereumClient *ethclient.Client, assetsContract *market.Market, freeverse *ecdsa.PrivateKey) (*Processor, error) {
	return &Processor{db, ethereumClient, assetsContract, freeverse, NewSigner(assetsContract)}, nil
}

func (b *Processor) Process() error {
	log.Info("Processing")

	openedAuctions, err := b.db.GetOpenAuctions()
	if err != nil {
		return err
	}

	for _, auction := range openedAuctions {
		bids, err := b.db.GetBidsOfAuction(auction.UUID)
		if err != nil {
			return err
		}

		machine, err := auctionmachine.New(
			auction,
			bids,
			b.assets,
			b.freeverse,
		)
		if err != nil {
			return err
		}

		err = machine.Process()
		if err != nil {
			return err
		}

		// update auction state if changed
		newState := machine.Auction.State
		if newState != auction.State {
			log.Infof("Auction %v: %v -> %v", auction.UUID, auction.State, newState)
			switch newState {
			case storage.AUCTION_PAYING:
				err = b.db.UpdateAuctionPaymentUrl(auction.UUID, "https://www.freeverse.io")
				if err != nil {
					return err
				}
				bid := bids[0]
				err = b.db.UpdateBidState(bid.Auction, bid.ExtraPrice, storage.BID_PAYING)
				if err != nil {
					return err
				}
				err = b.db.UpdateBidPaymentUrl(bid.Auction, bid.ExtraPrice, "http://ninjaflex.com/")
				if err != nil {
					return err
				}
				break
			}

			err = b.db.UpdateAuctionState(auction.UUID, newState)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Processor) FreezePlayer(Auction storage.Auction) error {
	sellerHiddenPrice, err := b.signer.HashPrivateMsg(
		Auction.CurrencyID,
		Auction.Price,
		Auction.Rnd,
	)
	if err != nil {
		return err
	}
	var sigs [3][32]byte
	var vs uint8
	sigs[0], err = b.signer.HashSellMessage(
		Auction.CurrencyID,
		Auction.Price,
		Auction.Rnd,
		Auction.ValidUntil,
		Auction.PlayerID,
	)
	if err != nil {
		return err
	}
	sigs[1], sigs[2], vs, err = b.signer.RSV(Auction.Signature)
	if err != nil {
		log.Error(err)
	}
	tx, err := b.assets.FreezePlayer(
		bind.NewKeyedTransactor(b.freeverse),
		sellerHiddenPrice,
		Auction.ValidUntil,
		Auction.PlayerID,
		sigs,
		vs,
	)
	if err != nil {
		return err
	}
	err = b.waitReceipt(tx, 10)
	if err != nil {
		return err
	}
	return nil
}

func (b *Processor) waitReceipt(tx *types.Transaction, timeoutSec uint8) error {
	receiptTimeout := time.Second * time.Duration(timeoutSec)
	start := time.Now()
	ctx := context.TODO()
	var receipt *types.Receipt

	for receipt == nil && time.Now().Sub(start) < receiptTimeout {
		receipt, err := b.client.TransactionReceipt(ctx, tx.Hash())
		if err == nil && receipt != nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return errors.New("Timeout waiting for receipt")
}
