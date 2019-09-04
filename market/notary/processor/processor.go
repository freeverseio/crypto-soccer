package processor

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/market/notary/contracts/assets"
	"github.com/freeverseio/crypto-soccer/market/notary/storage"
	log "github.com/sirupsen/logrus"
)

type Processor struct {
	db        *storage.Storage
	client    *ethclient.Client
	assets    *assets.Assets
	freeverse *ecdsa.PrivateKey
}

func NewProcessor(db *storage.Storage, ethereumClient *ethclient.Client, assetsContract *assets.Assets, freeverse *ecdsa.PrivateKey) (*Processor, error) {
	return &Processor{db, ethereumClient, assetsContract, freeverse}, nil
}

func (b *Processor) Process() error {
	log.Info("Processing")

	orders, err := b.db.GetOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		log.Infof("[broker] player %v -> team %v", order.SellOrder.PlayerId, order.BuyOrder.TeamId)
		playerId := big.NewInt(int64(order.SellOrder.PlayerId))
		teamId := big.NewInt(int64(order.BuyOrder.TeamId))
		_, err = b.assets.TransferPlayer(bind.NewKeyedTransactor(b.freeverse), playerId, teamId)
		if err != nil {
			log.Error(err)
		}
		err = b.db.DeleteOrder(order.SellOrder.PlayerId)
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}
