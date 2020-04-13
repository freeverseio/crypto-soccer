package consumer

import (
	"github.com/freeverseio/crypto-soccer/go/notary/producer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	ch chan interface{}
}

func NewConsumer(
	ch chan interface{},
) *Consumer {
	return &Consumer{
		ch,
	}
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
		default:
			log.Errorf("unknown event: %v", event)
		}
	}
}
