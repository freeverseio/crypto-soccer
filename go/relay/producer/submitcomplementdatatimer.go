package producer

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type SubmitComplementData struct{}

func NewSubmitComplementDataTimer(c chan interface{}, duration time.Duration) {
	for {
		time.Sleep(duration)
		if c != nil {
			select {
			case c <- SubmitComplementData{}:
			default:
				log.Warning("SubmitComplementDataTime::Start, channel is full, discarding value")
			}
		}
	}
}
