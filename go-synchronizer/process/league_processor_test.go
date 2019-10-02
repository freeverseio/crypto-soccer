package process_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/updates"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/testutils"
)

func TestProcessInvalidTimezone(t *testing.T) {
	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ganache := testutils.NewGanache()
	ganache.DeployContracts(ganache.Owner)
	processor, err := process.NewLeagueProcessor(ganache.Engine, ganache.Leagues, sto)
	if err != nil {
		t.Fatal(err)
	}
	var event updates.UpdatesActionsSubmission
	event.TimeZone = 25
	err = processor.Process(event)
	if err == nil {
		t.Fatal("processing invalid timezone")
	}
}

// func TestProcess(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../sql/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	ganache := testutils.NewGanache()
// 	ganache.DeployContracts(ganache.Owner)

// 	timezoneIdx := uint8(1)
// 	sto.TimezoneCreate(storage.Timezone{timezoneIdx})
// 	countryIdx := uint32(0)
// 	sto.CountryCreate(storage.Country{timezoneIdx, countryIdx})
// 	leagueIdx := uint32(0)
// 	sto.LeagueCreate(storage.League{timezoneIdx, countryIdx, leagueIdx})

// 	processor, err := process.NewLeagueProcessor(ganache.Engine, ganache.Leagues, sto)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var event updates.UpdatesActionsSubmission
// 	event.Day = 1
// 	event.TimeZone = timezoneIdx
// 	event.TurnInDay = 1
// 	err = processor.Process(event)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
