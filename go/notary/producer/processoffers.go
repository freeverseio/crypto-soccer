package producer

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type ProcessOfferEvent struct{}

func NewProcessorOffer(ch chan interface{}, duration time.Duration) {
	if ch == nil {
		log.Error("Nil channel")
		return
	}

	for {
		time.Sleep(duration)
		select {
		case ch <- ProcessOfferEvent{}:
		default:
			log.Warning("channel is full")
		}
	}
}
