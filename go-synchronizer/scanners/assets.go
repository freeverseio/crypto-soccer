package scanners

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
)

func ScanTeamCreated(assetsContract *assets.Assets) ([]assets.AssetsTeamCreated, error) {
	iter, err := assetsContract.FilterTeamCreated(&bind.FilterOpts{Start: 0})
	if err != nil {
		return nil, err
	}

	events := []assets.AssetsTeamCreated{}

	for iter.Next() {
		events = append(events, *(iter.Event))
	}
	return events, nil
}
