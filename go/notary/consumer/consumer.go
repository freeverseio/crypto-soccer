package consumer

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/producer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	ch chan interface{}
	db *sql.DB
}

func New(
	ch chan interface{},
	db *sql.DB,
) (*Consumer, error) {
	consumer := Consumer{}
	consumer.ch = ch
	consumer.db = db
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
		case producer.ProcessEvent:
			log.Info("Received ProcessEvent")
			// auctions, err := storage.GetPendingAuctions()
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// for _, auction := range auctions {
			// 	auctionMachine, err := auctionmachine.New(*auction, nil, nil, nil)
			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	if err := auctionMachine.Process(nil); err != nil {
			// 		log.Fatal(err)
			// 	}
			// }
		default:
			log.Errorf("unknown event: %v", event)
		}
	}
}
