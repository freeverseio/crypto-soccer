package router

import (
	"errors"
	"math/big"
)

func right(x *big.Int, n uint) *big.Int {
	return new(big.Int).Rsh(x, n)
}

func left(x *big.Int, n uint) *big.Int {
	return new(big.Int).Lsh(x, n)
}

func and(x *big.Int, n int64) *big.Int {
	return new(big.Int).And(x, big.NewInt(n))
}

func andBN(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).And(x, y)
}

func or(x *big.Int, n int64) *big.Int {
	return new(big.Int).Or(x, big.NewInt(n))
}

func orBN(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Or(x, y)
}

func lessThan(x *big.Int, n int64) bool {
	return x.Cmp(big.NewInt(n)) == -1
}

func largerThan(x *big.Int, n int64) bool {
	return x.Cmp(big.NewInt(n)) == 1
}

func equals(x *big.Int, n int64) bool {
	return x.Cmp(big.NewInt(n)) == 0
}

func twoToPow(n uint64) int64 {
	return 2 << (n - 1)
}

func encodeTZCountryAndValGo(timeZone uint8, countryIdxInTZ *big.Int, val *big.Int) (*big.Int, error) {
	if !(timeZone < 32) {
		return nil, errors.New("timezone out of bound")
	}
	if !lessThan(countryIdxInTZ, twoToPow(10)) {
		return nil, errors.New("countryIdxInTZ out of bound")
	}
	if !lessThan(val, twoToPow(28)) {
		return nil, errors.New("val out of bound")
	}
	encoded := orBN(left(big.NewInt(int64(timeZone)), 38), left(countryIdxInTZ, 28))
	return orBN(encoded, val), nil
}

func getCurrentShirtNumGo(playerState *big.Int) *big.Int {
	return and(right(playerState, 43), 31)
}

func getSkillGo(encodedSkills *big.Int, skillIdx uint8) *big.Int {
	return and(right(encodedSkills, uint(skillIdx)*20), 1048575) // 1048575 = 2**20 - 1
}

func getBirthDayGo(encodedSkills *big.Int) *big.Int {
	return and(right(encodedSkills, 100), 65535)
}

func getPlayerIdFromSkillsGo(encodedSkills *big.Int) *big.Int {
	if getIsSpecialGo(encodedSkills) {
		return encodedSkills
	}
	return getInternalPlayerIdGo(encodedSkills)
}

func getInternalPlayerIdGo(encodedSkills *big.Int) *big.Int {
	return and(right(encodedSkills, 129), 8796093022207) // 2**43 - 1 = 8796093022207
}

func getPotentialGo(encodedSkills *big.Int) *big.Int {
	return and(right(encodedSkills, 116), 15)
}

func getForwardnessGo(encodedSkills *big.Int) *big.Int {
	return and(right(encodedSkills, 120), 7)
}

func getLeftishnessGo(encodedSkills *big.Int) *big.Int {
	return and(right(encodedSkills, 123), 7)
}

func getAggressivenessGo(encodedSkills *big.Int) *big.Int {
	return and(right(encodedSkills, 126), 7)
}

func getAlignedEndOfFirstHalfGo(encodedSkills *big.Int) bool {
	return equals(and(right(encodedSkills, 172), 1), 1)
}

func getRedCardLastGameGo(encodedSkills *big.Int) bool {
	return equals(and(right(encodedSkills, 173), 1), 1)
}

func getGamesNonStoppingGo(encodedSkills *big.Int) uint8 {
	val := and(right(encodedSkills, 174), 7)
	return uint8(val.Uint64())
}

func getInjuryWeeksLeftGo(encodedSkills *big.Int) uint8 {
	val := and(right(encodedSkills, 177), 7)
	return uint8(val.Uint64())
}

func getSubstitutedFirstHalfGo(encodedSkills *big.Int) bool {
	return equals(and(right(encodedSkills, 180), 1), 1)
}

func getSumOfSkillsGo(encodedSkills *big.Int) *big.Int {
	return and(right(encodedSkills, 181), 524287) // 2**19-1
}

func getIsSpecialGo(encodedSkills *big.Int) bool {
	return equals(and(right(encodedSkills, 204), 1), 1)
}

func getGenerationGo(encodedSkills *big.Int) *big.Int {
	return and(right(encodedSkills, 205), 255)
}

func getOutOfGameFirstHalfGo(encodedSkills *big.Int) bool {
	return equals(and(right(encodedSkills, 213), 1), 1)
}

func getYellowCardFirstHalfGo(encodedSkills *big.Int) bool {
	return equals(and(right(encodedSkills, 214), 1), 1)
}

// MATCH LOG functions:

func getAssisterGo(log *big.Int, pos uint8) uint8 {
	val := and(right(log, (4+4*uint(pos))), 15)
	return uint8(val.Uint64())
}

func getShooterGo(log *big.Int, pos uint8) uint8 {
	val := and(right(log, (52+4*uint(pos))), 15)
	return uint8(val.Uint64())
}

func getForwardPosGo(log *big.Int, pos uint8) uint8 {
	val := and(right(log, (100+2*uint(pos))), 3)
	return uint8(val.Uint64())
}

func getPenaltyGo(log *big.Int, pos uint8) bool {
	return equals(and(right(log, 124+uint(pos)), 1), 1)
}

func getIsHomeStadiumGo(log *big.Int) bool {
	return equals(and(right(log, 248), 1), 1)
}

/// recall that 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
func getHalfTimeSubsGo(log *big.Int, pos uint8) uint8 {
	val := and(right(log, (179+5*uint(pos))), 31)
	return uint8(val.Uint64())
}

func getNGKAndDefsGo(log *big.Int, is2ndHalf bool) uint8 {
	offset := uint(194)
	if is2ndHalf {
		offset += 4
	}
	val := and(right(log, offset), 15)
	return uint8(val.Uint64())
}

func getNTotGo(log *big.Int, is2ndHalf bool) uint8 {
	offset := uint(202)
	if is2ndHalf {
		offset += 4
	}
	val := and(right(log, offset), 15)
	return uint8(val.Uint64())
}

func getWinnerGo(log *big.Int) uint8 {
	val := and(right(log, 210), 3)
	return uint8(val.Uint64())
}

func getTeamSumSkillsGo(log *big.Int) *big.Int {
	return and(right(log, 212), 16777215) // 2^24 - 1
}

func addTrainingPointsGo(log *big.Int, points *big.Int) *big.Int {
	return orBN(log, left(points, 236))
}

func getTrainingPointsGo(log *big.Int) uint16 {
	val := and(right(log, 236), 4095) // 2^12-1
	return uint16(val.Uint64())
}

func getNGoalsGo(log *big.Int) uint8 {
	val := and(log, 15)
	return uint8(val.Uint64())
}

func getOutOfGamePlayerGo(log *big.Int, is2ndHalf bool) *big.Int {
	var offset uint
	if is2ndHalf {
		offset = 141
	} else {
		offset = 131
	}
	return and(right(log, offset), 15)
}

func getOutOfGameTypeGo(log *big.Int, is2ndHalf bool) *big.Int {
	var offset uint
	if is2ndHalf {
		offset = 141
	} else {
		offset = 131
	}
	return and(right(log, offset+4), 3)
}

func getOutOfGameRoundGo(log *big.Int, is2ndHalf bool) *big.Int {
	var offset uint
	if is2ndHalf {
		offset = 141
	} else {
		offset = 131
	}
	return and(right(log, offset+6), 15)
}

func getYellowCardGo(log *big.Int, posInHaf uint8, is2ndHalf bool) uint8 {
	offset := uint(posInHaf) * 4
	if is2ndHalf {
		offset += 159
	} else {
		offset += 151
	}
	val := and(right(log, offset), 15)
	return uint8(val.Uint64())
}

func getInGameSubsHappenedGo(log *big.Int, posInHalf uint8, is2ndHalf bool) uint8 {
	offset := 167 + 2*uint(posInHalf)
	if is2ndHalf {
		offset += 6
	}
	val := and(right(log, offset), 3)
	return uint8(val.Uint64())
}

// TACTICS

func getTacticsIdGo(tactics *big.Int) uint8 {
	val := and(tactics, 63)
	return uint8(val.Uint64())
}

func getExtraAttackGo(tactics *big.Int, p uint8) bool {
	return equals(and(right(tactics, 6+uint(p)), 1), 1)
}

func getFullExtraAttackGo(tactics *big.Int) [10]bool {
	var extraAttack [10]bool
	for p := uint8(0); p < 10; p++ {
		extraAttack[p] = getExtraAttackGo(tactics, p)
	}
	return extraAttack
}

func getSubstitutionGo(tactics *big.Int, p uint8) uint8 {
	val := right(and(tactics, 86+4*int64(p)), 15)
	return uint8(val.Uint64())
}

func getSubsRoundGo(tactics *big.Int, p uint8) uint8 {
	val := right(and(tactics, 98+4*int64(p)), 15)
	return uint8(val.Uint64())
}

func getLinedUpGo(tactics *big.Int, p uint8) uint8 {
	val := right(and(tactics, 16+5*int64(p)), 31)
	return uint8(val.Uint64())
}

func getFullLineUpGo(tactics *big.Int) [14]uint8 {
	var lineup [14]uint8
	for p := uint8(0); p < 14; p++ {
		lineup[p] = getLinedUpGo(tactics, p)
	}
	return lineup
}
