package process

import (
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

func Process(assetsContract Assets, sto storage.Storage) {
	log.Debug("Process: called")

	log.Debug("Process: scanning the blockchain")
	countTeams, err := assetsContract.CountTeams(nil)
	if err != nil {
		log.Fatalf("Failed to retrieve token name: %v", err)
	}

	log.Debug("Process: act on local storage")
	for i := 0; i < int(countTeams.Uint64()); i++ {
		err = sto.TeamAdd(uint64(i), "name")
		if err != nil {
			log.Fatal(err)
		}
	}
}
