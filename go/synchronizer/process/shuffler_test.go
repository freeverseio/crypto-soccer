package process_test

import (
	"fmt"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"gotest.tools/assert"
)

func TestTimezoneToReshuffle(t *testing.T) {
	cases := []struct {
		TimezoneIdx uint8
		Day         uint8
		TurnInDay   uint8
		Result      uint8
	}{
		{0, 0, 0, 0},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			result := process.TimezoneToReshuffle(
				tc.TimezoneIdx,
				tc.Day,
				tc.TurnInDay,
			)
			assert.Equal(t, tc.Result, result)
		})
	}
}
