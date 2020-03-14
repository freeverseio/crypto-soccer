package matchevents_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/matchevents"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func TestMatchEventsNewMatchEvents(t *testing.T) {
	verseSeed := [32]byte{0x01}
	homeTeamID := "1"
	visitorTeamID := "2"
	homeTactic, _ := new(big.Int).SetString("340596594427581673436941882753025", 10)
	visitorTactic, _ := new(big.Int).SetString("340596594427581673436941882753025", 10)
	logsAndEvents := []*big.Int{}
	is2ndHalf := false
	events, err := matchevents.NewMatchEvents(
		*bc.Contracts,
		verseSeed,
		homeTeamID,
		visitorTeamID,
		homeTactic,
		visitorTactic,
		logsAndEvents,
		is2ndHalf,
	)
	assert.NilError(t, err)
	golden.Assert(t, events.DumpState(), t.Name()+".golden")
}
