package processor

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/market/go-broker/storage"
	log "github.com/sirupsen/logrus"
)

type Processor struct {
	db     *storage.Storage
	client *ethclient.Client
}

func NewProcessor(db *storage.Storage, ethereumClient string) *Processor {
	log.Info("Dial the Ethereum client: ", ethereumClient)
	client, err := ethclient.Dial(ethereumClient)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	return &Processor{db, client}
}

func (b *Processor) Process() {
	log.Info("Processing")

}
