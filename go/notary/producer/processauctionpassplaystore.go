package producer

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type AuctionPassPlaystoreOrderEvent struct{}

func NewAuctionPassPlaystoreOrderEventProcessor(ch chan interface{}, duration time.Duration) {
	if ch == nil {
		log.Error("Nil channer")
		return
	}

	for {
		time.Sleep(duration)
		select {
		case ch <- AuctionPassPlaystoreOrderEvent{}:
		default:
			log.Warning("channel is full")
		}
	}
}
