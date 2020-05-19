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
	ch                chan interface{}
	db                *sql.DB
	contracts         contracts.Contracts
	pvc               *ecdsa.PrivateKey
	market            marketpay.IMarketPay
	googleCredentials []byte
	iapTestOn         bool
}

func New(
	ch chan interface{},
	market marketpay.IMarketPay,
	db *sql.DB,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	googleCredentials []byte,
	iapTestOn bool,
) (*Consumer, error) {
	consumer := Consumer{}
	consumer.ch = ch
	consumer.db = db
	consumer.contracts = contracts
	consumer.pvc = pvc
	consumer.market = market
	consumer.googleCredentials = googleCredentials
	consumer.iapTestOn = iapTestOn
	return &consumer, nil
}

func (b *Consumer) Consume(event interface{}) error {
	switch in := event.(type) {
	case input.CreateAuctionInput:
		log.Info("Received CreateAuctionInput")
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
		log.Info("Received CancelAuctionInput")
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
		log.Info("Received CreateBidInput")
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
		log.Info("Received ProcessEvent")
		tx, err := b.db.Begin()
		if err != nil {
			return err
		}
		if err := ProcessAuctions(
			b.market,
			tx,
			b.contracts,
			b.pvc,
		); err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	case producer.PlaystoreOrderEvent:
		log.Info("Received PlaystoreOrderEvent")
		tx, err := b.db.Begin()
		if err != nil {
			return err
		}
		if err := ProcessPlaystoreOrders(
			tx,
			b.contracts,
			b.pvc,
			b.googleCredentials,
			b.iapTestOn,
		); err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	case input.SubmitPlayStorePlayerPurchaseInput:
		log.Info("Received SubmitPlayStorePlayerPurchaseInput")
		tx, err := b.db.Begin()
		if err != nil {
			return err
		}
		if err := SubmitPlayStorePlayerPurchase(
			tx,
			in,
		); err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown event: %+v", event)
	}
	return nil
}

func (b *Consumer) Start() {
	for {
		event := <-b.ch
		if err := b.Consume(event); err != nil {
			log.Error(err)
		}
	}
}
