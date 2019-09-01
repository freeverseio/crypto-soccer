package processor

import (
	"github.com/freeverseio/crypto-soccer/market/go-broker/storage"
	log "github.com/sirupsen/logrus"
)

type Processor struct {
	db *storage.Storage
}

func NewProcessor(db *storage.Storage) *Processor {
	return &Processor{db}
}

func (b *Processor) Process() {
	log.Info("Processing")

}
