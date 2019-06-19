package assets

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func scanTeamCreated(assets *Assets) ([]AssetsTeamCreated, error) {
	iter, err := assets.FilterTeamCreated(&bind.FilterOpts{Start: 0})
	if err != nil {
		return nil, err
	}

	events := []AssetsTeamCreated{}

	for iter.Next() {
		events = append(events, *(iter.Event))
	}
	return events, nil
}
