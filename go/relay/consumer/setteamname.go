package consumer

import (
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/storage"
	log "github.com/sirupsen/logrus"
)

func SetTeamName(teamStorageService storage.TeamStorageService, ev input.SetTeamNameInput) error {
	log.Infof("[relay|consumer] set team %v name %v", string(ev.TeamId), ev.Name)
	return teamStorageService.UpdateName(string(ev.TeamId), ev.Name)
}
