package router

import (
	"errors"
	"fmt"
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

func notBN(x *big.Int) *big.Int {
	return new(big.Int).Not(x)
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

func setValAtPos(serialized *big.Int, val uint, pos uint, mask *big.Int) *big.Int {
	bigVal := big.NewInt(int64(val))
	return orBN(andBN(serialized, notBN(left(mask, pos))), left(bigVal, pos))
}

func getValAtPos(serialized *big.Int, pos uint, mask *big.Int) int64 {
	return andBN(right(serialized, pos), mask).Int64()
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
	val := and(right(tactics, 86+4*uint(p)), 15)
	return uint8(val.Uint64())
}

func getSubsRoundGo(tactics *big.Int, p uint8) uint8 {
	val := and(right(tactics, 98+4*uint(p)), 15)
	return uint8(val.Uint64())
}

func getLinedUpGo(tactics *big.Int, p uint8) uint8 {
	val := and(right(tactics, 16+5*uint(p)), 31)
	return uint8(val.Uint64())
}

func getFullLineUpGo(tactics *big.Int) [14]uint8 {
	var lineup [14]uint8
	for p := uint8(0); p < 14; p++ {
		lineup[p] = getLinedUpGo(tactics, p)
	}
	return lineup
}

// MATCH EVENTS

func SetTeamThatAttacks(eventsLog *big.Int, round uint, teamThatAttacks uint) (*big.Int, error) {
	N_ROUNDS := uint(12)
	if !(round < N_ROUNDS) {
		return eventsLog, errors.New("round is too large")
	}
	if !(teamThatAttacks < 2) {
		return eventsLog, errors.New("teamThatAttacks is too large")
	}
	return setValAtPos(eventsLog, teamThatAttacks, 11*round, big.NewInt(1)), nil
}

func GetTeamThatAttacks(eventsLog *big.Int, round uint) (uint, error) {
	N_ROUNDS := uint(12)
	if !(round < N_ROUNDS) {
		return 0, errors.New("round is too large")
	}
	return uint(getValAtPos(eventsLog, 11*round, big.NewInt(1))), nil
}

func SetShooter(eventsLog *big.Int, round uint, player uint) (*big.Int, error) {
	N_ROUNDS := uint(12)
	MAX_PLAYER := uint(15)
	META_PLAYER := uint(100)

	if !(round < N_ROUNDS) {
		return eventsLog, errors.New("round is too large")
	}
	if (player > MAX_PLAYER) && (player != META_PLAYER) {
		return eventsLog, errors.New("player is too large: " + fmt.Sprint(player))
	}
	return setValAtPos(eventsLog, player, 11*round+1, big.NewInt(15)), nil
}

func GetShooter(eventsLog *big.Int, round uint) (uint, error) {
	N_ROUNDS := uint(12)
	if !(round < N_ROUNDS) {
		return 0, errors.New("round is too large")
	}
	return uint(getValAtPos(eventsLog, 11*round+1, big.NewInt(15))), nil
}

func SetAssister(eventsLog *big.Int, round uint, player uint) (*big.Int, error) {
	N_ROUNDS := uint(12)
	MAX_PLAYER := uint(15)
	META_PLAYER := uint(100)

	if !(round < N_ROUNDS) {
		return eventsLog, errors.New("round is too large")
	}
	if (player > MAX_PLAYER) && (player != META_PLAYER) {
		return eventsLog, errors.New("player is too large: " + fmt.Sprint(player))
	}
	return setValAtPos(eventsLog, player, 11*round+6, big.NewInt(15)), nil
}

func GetAssister(eventsLog *big.Int, round uint) (uint, error) {
	N_ROUNDS := uint(12)
	if !(round < N_ROUNDS) {
		return 0, errors.New("round is too large")
	}
	return uint(getValAtPos(eventsLog, 11*round+6, big.NewInt(15))), nil
}

func SetIsGoal(eventsLog *big.Int, round uint, isGoal bool) (*big.Int, error) {
	N_ROUNDS := uint(12)
	if !(round < N_ROUNDS) {
		return eventsLog, errors.New("round is too large")
	}
	val := uint(0)
	if isGoal {
		val = 1
	}
	return setValAtPos(eventsLog, val, 11*round+5, big.NewInt(1)), nil
}

func GetIsGoal(eventsLog *big.Int, round uint) (bool, error) {
	N_ROUNDS := uint(12)
	if !(round < N_ROUNDS) {
		return false, errors.New("round is too large")
	}
	return getValAtPos(eventsLog, 11*round+5, big.NewInt(1)) == int64(1), nil
}

func SetManagesToShoot(eventsLog *big.Int, round uint, managesToShoot bool) (*big.Int, error) {
	N_ROUNDS := uint(12)
	if !(round < N_ROUNDS) {
		return eventsLog, errors.New("round is too large")
	}
	val := uint(0)
	if managesToShoot {
		val = 1
	}
	return setValAtPos(eventsLog, val, 11*round+10, big.NewInt(1)), nil
}

func GetManagesToShoot(eventsLog *big.Int, round uint) (bool, error) {
	N_ROUNDS := uint(12)
	if !(round < N_ROUNDS) {
		return false, errors.New("round is too large")
	}
	return getValAtPos(eventsLog, 11*round+10, big.NewInt(1)) == int64(1), nil
}

func EncodeMatchEvents(
	teamThatAttacks []uint,
	shooter []uint,
	assister []uint,
	isGoal []bool,
	managesToShoot []bool,
) (*big.Int, error) {
	eventsLog := big.NewInt(0)
	var err error
	nRounds := len(teamThatAttacks)
	if len(shooter) != nRounds || len(assister) != nRounds || len(isGoal) != nRounds || len(managesToShoot) != nRounds {
		return eventsLog, errors.New("inputs to EncodeMatchEvents have different size")
	}
	for r := uint(0); r < uint(len(teamThatAttacks)); r++ {
		eventsLog, err = SetTeamThatAttacks(eventsLog, r, teamThatAttacks[r])
		if err != nil {
			return eventsLog, err
		}
		eventsLog, err = SetShooter(eventsLog, r, shooter[r])
		if err != nil {
			return eventsLog, err
		}
		eventsLog, err = SetAssister(eventsLog, r, assister[r])
		if err != nil {
			return eventsLog, err
		}
		eventsLog, err = SetIsGoal(eventsLog, r, isGoal[r])
		if err != nil {
			return eventsLog, err
		}
		eventsLog, err = SetManagesToShoot(eventsLog, r, managesToShoot[r])
		if err != nil {
			return eventsLog, err
		}
	}
	return eventsLog, nil
}

func DecodeMatchEvents(eventsLog *big.Int, nRounds uint) (
	[]uint,
	[]bool,
	[]uint,
	[]bool,
	[]uint,
	error,
) {
	var teamThatAttacks []uint
	var shooter []uint
	var assister []uint
	var isGoal []bool
	var managesToShoot []bool
	var err error
	var in uint
	var bo bool
	for r := uint(0); r < nRounds; r++ {
		in, err = GetTeamThatAttacks(eventsLog, r)
		if err != nil {
			return teamThatAttacks, managesToShoot, shooter, isGoal, assister, err
		}
		teamThatAttacks = append(teamThatAttacks, in)

		in, err = GetShooter(eventsLog, r)
		if err != nil {
			return teamThatAttacks, managesToShoot, shooter, isGoal, assister, err
		}
		shooter = append(shooter, in)

		in, err = GetAssister(eventsLog, r)
		if err != nil {
			return teamThatAttacks, managesToShoot, shooter, isGoal, assister, err
		}
		assister = append(assister, in)

		bo, err = GetIsGoal(eventsLog, r)
		if err != nil {
			return teamThatAttacks, managesToShoot, shooter, isGoal, assister, err
		}
		isGoal = append(isGoal, bo)

		bo, err = GetManagesToShoot(eventsLog, r)
		if err != nil {
			return teamThatAttacks, managesToShoot, shooter, isGoal, assister, err
		}
		managesToShoot = append(managesToShoot, bo)
	}
	return teamThatAttacks, managesToShoot, shooter, isGoal, assister, nil
}

func SerializeEventsFromPlayHalf(
	matchEvents []*big.Int,
) (*big.Int, error) {
	// // In the future, when we enable serialization directly from Solidity, the following "if" will be executed
	isAlreadyCompressed := equals(matchEvents[0], 2)
	if isAlreadyCompressed {
		return matchEvents[1], nil
	}

	eventsLog := big.NewInt(0)
	var err error
	if !(len(matchEvents)%5 == 0) {
		return eventsLog, errors.New("the length of matchEvents should be a multiple of 5")
	}
	nRounds := uint(len(matchEvents) / 5)
	for round := uint(0); round < nRounds; round++ {
		eventsLog, err = SetTeamThatAttacks(eventsLog, round, uint(matchEvents[5*round].Int64()))
		if err != nil {
			return eventsLog, err
		}
		eventsLog, err = SetManagesToShoot(eventsLog, round, equals(matchEvents[1+5*round], 1))
		if err != nil {
			return eventsLog, err
		}
		eventsLog, err = SetShooter(eventsLog, round, uint(matchEvents[2+5*round].Int64()))
		if err != nil {
			return eventsLog, err
		}
		eventsLog, err = SetIsGoal(eventsLog, round, equals(matchEvents[3+5*round], 1))
		if err != nil {
			return eventsLog, err
		}
		eventsLog, err = SetAssister(eventsLog, round, uint(matchEvents[4+5*round].Int64()))
		if err != nil {
			return eventsLog, err
		}

	}
	return eventsLog, nil
}
