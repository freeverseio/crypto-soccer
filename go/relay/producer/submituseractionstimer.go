package producer

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type SubmitActionsEvent struct{}

func NewSubmitUserActionsTimer(c chan interface{}, duration time.Duration) {
	for {
		time.Sleep(duration)
		if c != nil {
			select {
			case c <- SubmitActionsEvent{}:
			default:
				log.Warning("SubmitTimer::Start, channel is full, discarding value")
			}
		}
	}
}
