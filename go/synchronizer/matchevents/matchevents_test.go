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
	logsAndEvents := []*big.Int{}
	is2ndHalf := false
	_, err := matchevents.NewMatchEvents(
		*bc.Contracts,
		verseSeed,
		homeTeamID,
		visitorTeamID,
		homeTactic,
		visitorTactic,
		logsAndEvents,
		is2ndHalf,
	)
	assert.Error(t, err, "logAndEvents len < 2")
}

func TestMatchEventsNewMatchEventsOnlyMatchLogs(t *testing.T) {
	verseSeed := [32]byte{0x01}
	homeTeamID := "1"
	visitorTeamID := "2"
	homeTactic, _ := new(big.Int).SetString("340596594427581673436941882753025", 10)
	visitorTactic, _ := new(big.Int).SetString("340596594427581673436941882753025", 10)
	logsAndEvents := []*big.Int{}
	logsAndEvents[0], _ = new(big.Int).SetString("1809252841225230840719990802576567867449345403164777349740842651283017957376", 10)
	logsAndEvents[1], _ = new(big.Int).SetString("1809252842383666049074119856253627407002480579871841461333007538180245710739", 10)
	is2ndHalf := false
	_, err := matchevents.NewMatchEvents(
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
}
