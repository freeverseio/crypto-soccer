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
		{1, 0, 0, 2},
		{1, 1, 0, 0},
		{1, 0, 1, 0},
		{1, 13, 0, 0},
		{22, 13, 0, 0},
		{23, 12, 0, 0},
		{23, 13, 1, 0},
		{23, 13, 0, 24},
		{24, 13, 0, 1},
		{23, 13, 1, 0},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			result, err := process.TimezoneToReshuffle(
				tc.TimezoneIdx,
				tc.Day,
				tc.TurnInDay,
			)
			assert.Equal(t, err, nil)
			assert.Equal(t, tc.Result, result)
		})
	}
}

func TestNonValidTZ(t *testing.T) {
	cases := []struct {
		TimezoneIdx uint8
		Day         uint8
		TurnInDay   uint8
		Result      uint8
	}{
		{0, 0, 0, 2},
		{25, 1, 0, 0},
		{27, 0, 1, 0},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			_, err := process.TimezoneToReshuffle(
				tc.TimezoneIdx,
				tc.Day,
				tc.TurnInDay,
			)
			assert.Equal(t, err == nil, false)
		})
	}
}
