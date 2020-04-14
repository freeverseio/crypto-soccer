package producer

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type ProcessEvent struct {
	c chan interface{}
}

func NewProcessor(c chan interface{}, duration time.Duration) {
	for {
		time.Sleep(duration)
		if c != nil {
			select {
			case c <- ProcessEvent{}:
			default:
				log.Warning("channel is full")
			}
		}
	}
}
