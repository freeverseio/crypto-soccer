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

func (b *Processor) processOrder(order storage.Order) error {
	log.Infof("[broker] player %v -> team %v", order.SellOrder.PlayerId, order.BuyOrder.TeamId)

	log.Infof("(1) generate hash private msg")
	privHash, err := b.signer.HashPrivateMsg(
		order.SellOrder.CurrencyId,
		order.SellOrder.Price,
		order.SellOrder.Rnd,
	)
	if err != nil {
		return err
	}

	log.Infof("(2) generate hash sell message")
	var sigs [6][32]byte
	var vs [2]uint8
	sigs[0], err = b.signer.HashSellMessage(
		order.SellOrder.CurrencyId,
		order.SellOrder.Price,
		order.SellOrder.Rnd,
		order.SellOrder.ValidUntil,
		order.SellOrder.PlayerId,
		order.SellOrder.TypeOfTx,
	)
	if err != nil {
		return err
	}
	sigs[1], sigs[2], vs[0], err = b.signer.RSV(order.SellOrder.Signature)
	if err != nil {
		log.Error(err)
	}
	log.Infof("(3) generate hash buy message")
	sigs[3], err = b.signer.HashBuyMessage(
		order.SellOrder.CurrencyId,
		order.SellOrder.Price,
		order.SellOrder.Rnd,
		order.SellOrder.ValidUntil,
		order.SellOrder.PlayerId,
		order.SellOrder.TypeOfTx,
		order.BuyOrder.TeamId,
	)
	if err != nil {
		return err
	}
	sigs[4], sigs[5], vs[1], err = b.signer.RSV(order.BuyOrder.Signature)
	if err != nil {
		return err
	}

	log.Infof("(4) freeze player")
	tx, err := b.assets.FreezePlayer(
		bind.NewKeyedTransactor(b.freeverse),
		privHash,
		order.SellOrder.ValidUntil,
		order.SellOrder.PlayerId,
		order.SellOrder.TypeOfTx,
		order.BuyOrder.TeamId,
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
	log.Infof("(5) complete freeze")
	tx, err = b.assets.CompleteFreeze(
		bind.NewKeyedTransactor(b.freeverse),
		order.SellOrder.PlayerId,
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

func (b *Processor) Process() error {
	log.Info("Processing")

	orders, err := b.db.GetOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		err = b.processOrder(order)
		if err != nil {
			log.Error(err)
		}

		log.Infof("(CLEANING) delete order")
		err = b.db.DeleteOrder(order.SellOrder.PlayerId)
		if err != nil {
			log.Error(err)
		}
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
