package consumer

import (
	"crypto/ecdsa"
	"database/sql"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/producer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	ch        chan interface{}
	db        *sql.DB
	contracts contracts.Contracts
	pvc       *ecdsa.PrivateKey
	market    marketpay.IMarketPay
}

func New(
	ch chan interface{},
	market marketpay.IMarketPay,
	db *sql.DB,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
) (*Consumer, error) {
	consumer := Consumer{}
	consumer.ch = ch
	consumer.db = db
	consumer.contracts = contracts
	consumer.pvc = pvc
	consumer.market = market
	return &consumer, nil
}

func (b *Consumer) Consume(event interface{}) error {
	switch in := event.(type) {
	case input.CreateAuctionInput:
		log.Debug("Received CreateAuctionInput")
		tx, err := b.db.Begin()
		if err != nil {
			return err
		}
		if err := CreateAuction(tx, in); err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	case input.CancelAuctionInput:
		log.Debug("Received CancelAuctionInput")
		tx, err := b.db.Begin()
		if err != nil {
			return err
		}
		if err := CancelAuction(tx, in); err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	case input.CreateBidInput:
		log.Debug("Received CreateBidInput")
		tx, err := b.db.Begin()
		if err != nil {
			return err
		}
		if err := CreateBid(tx, in); err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	case producer.ProcessEvent:
		log.Debug("Received ProcessEvent")
		tx, err := b.db.Begin()
		if err != nil {
			return err
		}
		if err := ProcessAuctions(b.market, tx, b.contracts, b.pvc); err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	// case input.SubmitPlayStorePlayerPurchaseInput:
	// 	log.Debug("Received SubmitPlayStorePlayerPurchaseInput")
	default:
		return fmt.Errorf("unknown event: %+v", event)
	}
	return nil
}

func (b *Consumer) Start() {
	for {
		event := <-b.ch
		b.Consume(event)
	}
}
