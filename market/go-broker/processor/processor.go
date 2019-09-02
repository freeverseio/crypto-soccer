package processor

import (
	"math/big"

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

func (b *Processor) Process() error {
	log.Info("Processing")

	orders, err := b.db.GetOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		log.Infof("[broker] transfering player: %v", order.SellOrder.PlayerId)
		owner, err := b.assets.GetPlayerOwner(nil, big.NewInt(int64(order.SellOrder.PlayerId)))
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("owner : %v", owner)
	}

	return nil
}
