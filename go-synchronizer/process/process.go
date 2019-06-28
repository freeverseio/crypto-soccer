package process

import (
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/scanners"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

func Process(assetsContract *assets.Assets, sto *storage.Storage) error {
	log.Debug("Process: called")

	log.Debug("Process: scanning the blockchain")
	events, err := scanners.ScanTeamCreated(assetsContract, nil)
	if err != nil {
		return err
	}

	log.Debug("Process: act on local storage")
	for i := 0; i < len(events); i++ {
		event := events[i]
		err = sto.TeamAdd(event.Id.Uint64(), event.Name)
		if err != nil {
			return err
		}
	}

	return nil
}
