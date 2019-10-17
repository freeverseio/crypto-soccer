package processor

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/market/notary/contracts/market"
	"github.com/freeverseio/crypto-soccer/market/notary/storage"
	log "github.com/sirupsen/logrus"
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

// func (b *Processor) processOrder(order storage.Order) error {
// 	log.Infof("[broker] player %v -> team %v", order.Auction.PlayerId, order.Bid.TeamId)

// 	log.Infof("(1) generate hash private msg")
// 	sellerHiddenPrice, err := b.signer.HashPrivateMsg(
// 		order.Auction.CurrencyId,
// 		order.Auction.Price,
// 		order.Auction.Rnd,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	log.Infof("(2) generate hash sell message")
// 	var sigs [3][32]byte
// 	var vs uint8
// 	sigs[0], err = b.signer.HashSellMessage(
// 		order.Auction.CurrencyId,
// 		order.Auction.Price,
// 		order.Auction.Rnd,
// 		order.Auction.ValidUntil,
// 		order.Auction.PlayerId,
// 		order.Auction.TypeOfTx,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	sigs[1], sigs[2], vs, err = b.signer.RSV(order.Auction.Signature)
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	log.Infof("(3) generate hash buy message")
// 	_, err = b.signer.HashBuyMessage(
// 		order.Auction.CurrencyId,
// 		order.Auction.Price,
// 		order.Auction.Rnd,
// 		order.Auction.ValidUntil,
// 		order.Auction.PlayerId,
// 		order.Auction.TypeOfTx,
// 		order.Bid.TeamId,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	// sigs[4], sigs[5], vs[1], err = b.signer.RSV(order.Bid.Signature)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	log.Infof("(4) freeze player")
// 	tx, err := b.assets.FreezePlayer(
// 		bind.NewKeyedTransactor(b.freeverse),
// 		sellerHiddenPrice,
// 		order.Auction.ValidUntil,
// 		order.Auction.PlayerId,
// 		sigs,
// 		vs,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	err = b.waitReceipt(tx, 10)
// 	if err != nil {
// 		return err
// 	}
// 	log.Infof("(5) complete freeze")
// 	tx, err = b.assets.CompleteFreeze(
// 		bind.NewKeyedTransactor(b.freeverse),
// 		order.Auction.PlayerId,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	err = b.waitReceipt(tx, 10)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (b *Processor) Process() error {
	log.Info("Processing")

	// I get all the orders
	orders, err := b.db.GetOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		playerID := order.Auction.PlayerID
		frozen, err := b.assets.IsPlayerFrozen(&bind.CallOpts{}, playerID)
		if err != nil {
			log.Error(err)
			continue
		}
		if frozen == false {
			err = b.FreezePlayer(order.Auction)
			if err != nil {
				log.Error(err)
				continue
			}
		}

		// err = b.processOrder(order)
		// if err != nil {
		// 	log.Error(err)
		// }

		// log.Infof("(CLEANING) delete order")
		// err = b.db.DeleteOrder(order.Auction.PlayerId)
		// if err != nil {
		// 	log.Error(err)
		// }
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
