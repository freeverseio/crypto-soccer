package proxy2contracts

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
