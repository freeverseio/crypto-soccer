package consumer

import (
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/storage"
)

func SetTeamManagerName(teamStorageService storage.TeamStorageService, ev input.SetTeamManagerNameInput) error {
	return teamStorageService.UpdateManagerName(string(ev.TeamId), ev.Name)
}
