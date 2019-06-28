package process

import (
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/scanners"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

func Process(assetsContract *assets.Assets, sto *storage.Storage) error {
	log.Trace("Process: called")

	log.Trace("Process: scanning the blockchain")
	events, err := scanners.ScanTeamCreated(assetsContract, nil)
	if err != nil {
		return err
	}

	log.Trace("Process: act on local storage")
	for i := 0; i < len(events); i++ {
		event := events[i]
		name, err := assetsContract.GetTeamName(nil, event.Id)
		if err != nil {
			return err
		}
		err = sto.TeamAdd(event.Id.Uint64(), name)
		if err != nil {
			return err
		}
		log.Debugf("Team Created: id = %v, name = %v", event.Id.String(), name)
	}

	return nil
}
