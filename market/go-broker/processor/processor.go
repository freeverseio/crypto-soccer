package processor

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/market/go-broker/contracts/assets"
	"github.com/freeverseio/crypto-soccer/market/go-broker/storage"
	log "github.com/sirupsen/logrus"
)

type Processor struct {
	db     *storage.Storage
	client *ethclient.Client
	assets *assets.Assets
}

func NewProcessor(db *storage.Storage, ethereumClient string, assetsContractAddress string) (*Processor, error) {
	log.Info("Dial the Ethereum client: ", ethereumClient)
	client, err := ethclient.Dial(ethereumClient)
	if err != nil {
		return nil, err
	}
	log.Info("Creating Assets bindings to: ", assetsContractAddress)
	assetsContract, err := assets.NewAssets(common.HexToAddress(assetsContractAddress), client)
	if err != nil {
		return nil, err
	}
	return &Processor{db, client, assetsContract}, nil
}

func (b *Processor) Process() {
	log.Info("Processing")

}
