package consumer

import (
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/storage"
	log "github.com/sirupsen/logrus"
)

func SetTeamManagerName(teamStorageService storage.TeamStorageService, ev input.SetTeamManagerNameInput) error {
	log.Infof("[relay|consumer] set manager team %v name %v", string(ev.TeamId), ev.Name)
	return teamStorageService.UpdateManagerName(string(ev.TeamId), ev.Name)
}
