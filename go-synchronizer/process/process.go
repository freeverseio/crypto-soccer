package process

import (
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

func Process(assetsContract Assets, sto storage.Storage){
	log.Info("Process: called")

	countTeams, err := assetsContract.CountTeams(nil)
	if err != nil {
		log.Fatalf("Failed to retrieve token name: %v", err)
	}

	for i := 0; i < int(countTeams.Uint64()); i++ {
		sto.TeamAdd(uint64(i), "name")
	}

	// contractTeams, err := assetsContract.CountTeams(nil)
	// if err != nil {
	// 	log.Fatalf("Failed to retrieve token name: %v", err)
	// }
	// dbTeams, err := sto.TeamCount()
	// for i:=contractTeams;i<dbTeams;i++ {
	// 	teamName := fmt.Sptintf("team-%v",i)
	// 	assert(t,sto.TeamAdd(i,teamName))
	// }

}