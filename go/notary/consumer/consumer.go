package consumer

import (
	"crypto/ecdsa"
	"database/sql"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/names"
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
	namesdb           *names.Generator
	iapTestOn         bool
}

func New(
	ch chan interface{},
	market marketpay.IMarketPay,
	db *sql.DB,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	googleCredentials []byte,
	namesdb *names.Generator,
	iapTestOn bool,
) (*Consumer, error) {
	consumer := Consumer{}
	consumer.ch = ch
	consumer.db = db
	consumer.contracts = contracts
	consumer.pvc = pvc
	consumer.market = market
	consumer.googleCredentials = googleCredentials
	consumer.namesdb = namesdb
	consumer.iapTestOn = iapTestOn
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
		log.Info("[consumer] process auctions")
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
		log.Info("[consumer] process playstore events")
		tx, err := b.db.Begin()
		if err != nil {
			return err
		}
		if err := ProcessPlaystoreOrders(
			tx,
			b.contracts,
			b.pvc,
			b.googleCredentials,
			b.namesdb,
			b.iapTestOn,
		); err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	case input.SubmitPlayStorePlayerPurchaseInput:
		log.Debug("Received SubmitPlayStorePlayerPurchaseInput")
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
	case input.DismissPlayerInput:
		log.Debug("Received DismissPlayerInput")
		if err := DismissPlayer(
			b.contracts,
			b.pvc,
			in,
		); err != nil {
			return err
		}
	case input.CompletePlayerTransitInput:
		log.Debug("Received CompletePlayerTransit")
		if err := CompletePlayerTransit(b.contracts, b.pvc, in); err != nil {
			return err
		}
	case input.CreateOfferInput:
		log.Debug("Received CreateOfferInput")

		tx, err := b.db.Begin()
		if err != nil {
			return err
		}
		if err := CreateOffer(tx, in, b.contracts); err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	case input.AcceptOfferInput:
		log.Debug("Received CreateAuctionInput")
		tx, err := b.db.Begin()
		if err != nil {
			return err
		}
		if err := AcceptOffer(tx, in); err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	case input.CancelOfferInput:
		log.Debug("Received CancelOfferInput")
		tx, err := b.db.Begin()
		if err != nil {
			return err
		}
		if err := CancelOffer(tx, in); err != nil {
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
