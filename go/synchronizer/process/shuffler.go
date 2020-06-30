package process

import "errors"

// @return 0 in case of no timezone
func TimezoneToReshuffle(timezoneIdx uint8, day uint8, turnInDay uint8) (uint8, error) {
	// Reshuffle happens 1 hour before the start of a league
	// So for each timezone, on day 0, we reset timezone + 1.
	// This applies except for:
	// - the first timezone that played ever, i.e. tz = 24, which must be reset when tz = 23 plays day = 13
	if timezoneIdx < 1 || timezoneIdx > 24 {
		return 0, errors.New("input timezone must be in [1,24]")
	}
	timezoneToReshuffle := uint8(0)
	if (timezoneIdx == 23) && (day == 12) && (turnInDay == 0) {
		timezoneToReshuffle = 24
	}
	if (timezoneIdx == 24) && (day == 0) && (turnInDay == 0) {
		timezoneToReshuffle = 1
	}
	if (timezoneIdx < 23) && (day == 0) && (turnInDay == 0) {
		timezoneToReshuffle = timezoneIdx + 1
	}
	return timezoneToReshuffle, nil
}

func TimezoneToReshuffleOld(timezoneIdx uint8, day uint8, turnInDay uint8) (uint8, error) {
	// Reshuffle happens 1 hour before the start of a league
	// So for each timezone, on day 0, we reset timezone + 1.
	// This applies except for:
	// - the first timezone that played ever, i.e. tz = 24, which must be reset when tz = 23 plays day = 13
	if timezoneIdx < 1 || timezoneIdx > 24 {
		return 0, errors.New("input timezone must be in [1,24]")
	}
	timezoneToReshuffle := uint8(0)
	if (timezoneIdx == 23) && (day == 13) && (turnInDay == 0) {
		timezoneToReshuffle = 24
	}
	if (timezoneIdx == 24) && (day == 0) && (turnInDay == 0) {
		timezoneToReshuffle = 1
	}
	if (timezoneIdx < 23) && (day == 0) && (turnInDay == 0) {
		timezoneToReshuffle = timezoneIdx + 1
	}
	return timezoneToReshuffle, nil
}
