package producer

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type OfferEvent struct{}

func NewOfferEventProcessor(ch chan interface{}, duration time.Duration) {
	if ch == nil {
		log.Error("Nil channer")
		return
	}

	for {
		time.Sleep(duration)
		select {
		case ch <- OfferEvent{}:
		default:
			log.Warning("channel is full")
		}
	}
}
