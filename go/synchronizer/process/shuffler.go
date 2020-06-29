package process

// @return 0 in case of no timezone
func TimezoneToReshuffle(timezoneIdx uint8, day uint8, turnInDay uint8) uint8 {
	// Reshuffle happens 1 hour before the start of a league
	// So for each timezone, on day 0, we reset timezone + 1.
	// Except for the first timezone (tz=1), of course, which is reset the day before (day = 13)
	// during the match played by the last timezone (tz=24)
	timezoneToReshuffle := uint8(0)
	if (timezoneIdx == 24) && (day == 13) && (turnInDay == 0) {
		timezoneToReshuffle = 1
	}
	if (timezoneIdx < 24) && (day == 0) && (turnInDay == 0) {
		timezoneToReshuffle = timezoneIdx + 1
	}
	return timezoneToReshuffle
}
