package consumer

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/producer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	ch chan interface{}
	db *sql.DB
}

func New(
	ch chan interface{},
	db *sql.DB,
) *Consumer {
	consumer := Consumer{}
	consumer.ch = ch
	consumer.db = db
	return &consumer
}

func (b *Consumer) Start() {
	for {
		event := <-b.ch
		switch event.(type) {
		case gql.CreateAuctionInput:
			log.Debug("Received CreateAuctionInput")
		case gql.CancelAuctionInput:
			log.Debug("Received CancelAuctionInput")
		case gql.CreateBidInput:
			log.Debug("Received CreateBidInput")
		case producer.ProcessEvent:
			log.Debug("Received ProcessEvent")
			auctions, err := storage.GetPendingAuctions()
			if err != nil {
				log.Fatal(err)
			}
			for _, auction := range auctions {
				auctionMachine, err := auctionmachine.New(*auction, nil, nil, nil)
				if err != nil {
					log.Fatal(err)
				}
				if err := auctionMachine.Process(nil); err != nil {
					log.Fatal(err)
				}
			}
		default:
			log.Errorf("unknown event: %v", event)
		}
	}
}
