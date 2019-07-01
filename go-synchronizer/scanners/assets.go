package scanners

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
)

func ScanTeamCreated(assetsContract *assets.Assets, opts *bind.FilterOpts) ([]assets.AssetsTeamCreated, error) {
	if opts == nil {
		opts = &bind.FilterOpts{Start: 0}
	}
	iter, err := assetsContract.FilterTeamCreated(opts)
	if err != nil {
		return nil, err
	}

	events := []assets.AssetsTeamCreated{}

	for iter.Next() {
		events = append(events, *(iter.Event))
	}
	return events, nil
}
