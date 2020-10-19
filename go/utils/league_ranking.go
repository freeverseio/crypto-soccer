package utils

import "math/big"

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

func decodeTZCountryAndValGo(encoded *big.Int) (uint8, *big.Int, *big.Int) {
	return uint8(and(right(encoded, 38), 31).Int64()), and(right(encoded, 28), 1023), and(encoded, 268435455)
}

func getPerfPoints(leagueRanking uint8) uint64 {
	if leagueRanking == 0 {
		return 50
	} else if leagueRanking == 1 {
		return 42
	} else if leagueRanking == 2 {
		return 30
	} else if leagueRanking == 3 {
		return 25
	} else if leagueRanking == 4 {
		return 20
	} else if leagueRanking == 5 {
		return 15
	} else if leagueRanking == 6 {
		return 5
	} else {
		return 0
	}
}

func getSumOfTopPlayerSkills(skills []*big.Int) int64 {
	PLAYERS_PER_TEAM_MAX := uint8(25)
	var sortedSumSkills [25]int64
	for p := uint8(0); p < PLAYERS_PER_TEAM_MAX; p++ {
		if !equals(skills[p], 0) {
			sortedSumSkills[p] = getSumOfSkillsGo(skills[p]).Int64()
		}
	}
	sort25(&sortedSumSkills)
	teamSkills := int64(0)
	for p := uint8(0); p < 18; p++ {
		teamSkills += sortedSumSkills[p]
	}
	return teamSkills
}

func sort25(data *[25]int64) *[25]int64 {
	quickSort25(data, int(0), int(24))
	return data
}

func quickSort25(arr *[25]int64, left int, right int) {
	i := left
	j := right
	if i != j {
		pivot := arr[uint(left+(right-left)/2)]
		for i <= j {
			for arr[uint(i)] > pivot {
				i++
			}
			for pivot > arr[uint(j)] {
				j--
			}
			if i <= j {
				oldj := arr[uint(j)]
				oldi := arr[uint(i)]
				arr[uint(i)] = oldj
				arr[uint(j)] = oldi
				i++
				j--
			}
		}
		if left < j {
			quickSort25(arr, left, j)
		}
		if i < right {
			quickSort25(arr, i, right)
		}
	}
}

func getSumOfSkillsGo(encodedSkills *big.Int) *big.Int {
	return and(right(encodedSkills, 181), 524287) // 2**19-1
}

func computeTeamRankingPointsPure(
	skills [25]*big.Int,
	leagueRanking uint8,
	prevPerfPoints int64,
) (int64, int64) {
	WEIGHT_SKILLS := int64(100)
	INERTIA := int64(2)
	perfPointsThisLeague := getPerfPoints(leagueRanking)
	prevPerfPoints = (INERTIA*prevPerfPoints + (10-INERTIA)*perfPointsThisLeague)
	result := int64(getSumOfTopPlayerSkills(skills)) * (WEIGHT_SKILLS*10 + prevPerfPoints)
	return result, prevPerfPoints / 10
}
