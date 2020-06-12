package complementarydata

import "github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"

type ComplementaryData struct {
	SetTeamNameEvents        []input.SetTeamNameInput
	SetTeamManagerNameEvents []input.SetTeamManagerNameInput
}
