package processor

import (
	log "github.com/sirupsen/logrus"
)

type Processor struct {
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (b *Processor) Process() {
	log.Info("Processing")
}
