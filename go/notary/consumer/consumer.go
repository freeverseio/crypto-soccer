package consumer

import (
	"crypto/ecdsa"
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/producer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	ch        chan interface{}
	db        *sql.DB
	contracts contracts.Contracts
	pvc       *ecdsa.PrivateKey
}

func New(
	ch chan interface{},
	db *sql.DB,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
) (*Consumer, error) {
	consumer := Consumer{}
	consumer.ch = ch
	consumer.db = db
	consumer.contracts = contracts
	consumer.pvc = pvc
	return &consumer, nil
}

func (b *Consumer) Start() {
	for {
		event := <-b.ch
		switch in := event.(type) {
		case input.CreateAuctionInput:
			log.Debug("Received CreateAuctionInput")
			tx, err := b.db.Begin()
			if err != nil {
				log.Error(err)
				break
			}
			if err := CreateAuction(tx, in); err != nil {
				log.Error(err)
				tx.Rollback()
				break
			}
			if err = tx.Commit(); err != nil {
				log.Error(err)
			}
		case input.CancelAuctionInput:
			log.Debug("Received CancelAuctionInput")
			tx, err := b.db.Begin()
			if err != nil {
				log.Error(err)
				break
			}
			if err := CancelAuction(tx, in); err != nil {
				log.Error(err)
				tx.Rollback()
				break
			}
			if err = tx.Commit(); err != nil {
				log.Error(err)
			}
		case input.CreateBidInput:
			log.Debug("Received CreateBidInput")
			tx, err := b.db.Begin()
			if err != nil {
				log.Error(err)
				break
			}
			if err := CreateBid(tx, in); err != nil {
				log.Error(err)
				tx.Rollback()
				break
			}
			if err = tx.Commit(); err != nil {
				log.Error(err)
			}
		case producer.ProcessEvent:
			log.Debug("Received ProcessEvent")
			tx, err := b.db.Begin()
			if err != nil {
				log.Error(err)
				break
			}
			if err := ProcessAuctions(tx, b.contracts, b.pvc); err != nil {
				log.Fatal(err)
				tx.Rollback()
				break
			}
			if err = tx.Commit(); err != nil {
				log.Error(err)
			}
		default:
			log.Errorf("unknown event: %v", event)
		}
	}
}
