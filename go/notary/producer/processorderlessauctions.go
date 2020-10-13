package producer

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type ProcessOrderlessAuctionEvent struct{}

func NewProcessorOrderlessAuctions(ch chan interface{}, duration time.Duration) {
	if ch == nil {
		log.Error("Nil channer")
		return
	}

	for {
		time.Sleep(duration)
		select {
		case ch <- ProcessEvent{}:
		default:
			log.Warning("channel is full")
		}
	}
}
