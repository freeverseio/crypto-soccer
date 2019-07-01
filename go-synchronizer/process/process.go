package process

import (
	"context"
	//"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/scanners"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

func Process(assetsContract *assets.Assets, sto *storage.Storage, client *ethclient.Client) error {
	log.Info("Syncing ...")

	log.Trace("Process: scanning the blockchain")

	storedLastBlockNumber := big.NewInt(0)
	clientLastBlockNumber := big.NewInt(0)
	var events []assets.AssetsTeamCreated
	var err error

	if storedLastBlockNumber, err = sto.GetBlockNumber(); err != nil {
		return err
	}

	if client != nil {
		if header, err := client.HeaderByNumber(context.Background(), nil); err != nil {
			return err
		} else {
			clientLastBlockNumber = header.Number
		}
	} else {
		clientLastBlockNumber = big.NewInt(storedLastBlockNumber.Int64() + 1)
	}

	end := clientLastBlockNumber.Uint64()

	opts := &bind.FilterOpts{
		Start:   storedLastBlockNumber.Uint64(),
		End:     &end,
		Context: context.Background(),
	}

	if events, err = scanners.ScanTeamCreated(assetsContract, opts); err != nil {
		return err
	}

	//fmt.Println("Scanning from: ", storedLastBlockNumber.Uint64(), " to ", clientLastBlockNumber.Uint64())

	sto.SetBlockNumber(big.NewInt(clientLastBlockNumber.Int64() + 1))

	log.Trace("Process: act on local storage")
	for i := 0; i < len(events); i++ {
		event := events[i]
		name, err := assetsContract.GetTeamName(nil, event.Id)
		if err != nil {
			return err
		}
		//fmt.Println("TeamAdd", event.Id.Uint64(), name)
		err = sto.TeamAdd(event.Id.Uint64(), name)
		if err != nil {
			return err
		}
		log.Debugf("Team Created: id = %v, name = %v", event.Id.String(), name)
	}

	return nil
}
