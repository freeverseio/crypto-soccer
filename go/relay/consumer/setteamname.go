package consumer

import (
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/storage"
)

func SetTeamName(teamStorageService storage.TeamStorageService, ev input.SetTeamNameInput) error {
	return teamStorageService.UpdateName(string(ev.TeamId), ev.Name)
}
