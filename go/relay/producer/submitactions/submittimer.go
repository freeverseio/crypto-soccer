package submitactions

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type SubmitActionsEvent struct{}

func NewSubmitTimer(c chan interface{}, duration time.Duration) {
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
