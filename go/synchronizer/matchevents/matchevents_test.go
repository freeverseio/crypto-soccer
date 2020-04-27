package matchevents_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/matchevents"
	"gotest.tools/assert"
)

func TestMatchEventsNewMatchEventsNoMatchLogs(t *testing.T) {
	verseSeed := [32]byte{0x01}
	homeTeamID := "1"
	visitorTeamID := "2"
	homeTactic, _ := new(big.Int).SetString("340596594427581673436941882753025", 10)
	visitorTactic, _ := new(big.Int).SetString("340596594427581673436941882753025", 10)
	decodedHomeMatchLog := [15]uint32{}
	decodedVisitorMatchLog := [15]uint32{}
	logsAndEvents := []*big.Int{}
	is2ndHalf := false
	playerIds := [25]*big.Int{}
	for i := 0; i < 25; i++ {
		playerIds[i] = new(big.Int).SetUint64(21342314523)
	}
	_, err := matchevents.NewMatchEvents(
		*bc.Contracts,
		verseSeed,
		homeTeamID,
		visitorTeamID,
		playerIds,
		playerIds,
		homeTactic,
		visitorTactic,
		logsAndEvents,
		decodedHomeMatchLog,
		decodedVisitorMatchLog,
		is2ndHalf,
	)
	assert.Error(t, err, "logAndEvents len < 2")
}
