package worldplayer

import (
	"math/big"
	"strconv"
)

type WorldPlayersTier struct {
	LevelRange       [2]uint32
	PotentialWeights [10]uint32
	ProductId        string
	GoalKeepersCount uint8
	DefendersCount   uint8
	MidfieldersCount uint8
	AttackersCount   uint8
	Duration         int64
}

func addPlayerAtRandomFieldPos(tier WorldPlayersTier, seed string, randomPosPlayersCount int64) WorldPlayersTier {
	maxPos := uint64(4)
	for p := int64(0); p < randomPosPlayersCount; p++ {
		switch playerPos := generateRnd(big.NewInt(p), seed, maxPos); {
		case playerPos == 0:
			tier.GoalKeepersCount++
		case playerPos == 1:
			tier.DefendersCount++
		case playerPos == 2:
			tier.MidfieldersCount++
		case playerPos == 3:
			tier.AttackersCount++
		}
	}
	return tier
}

func generateBatchDistribution(teamId string, periodNumber int64) []WorldPlayersTier {
	var tiers []WorldPlayersTier

	seed := teamId + strconv.FormatUint(uint64(periodNumber), 10)

	// Tier1:
	// - has a fixed number of players, and fixed distribution of field position
	// - one reroll: 21 players. So if we want 1 player of 100% potential per week,
	// we need 1 in 21*2*7 = 294
	tier := WorldPlayersTier{
		LevelRange:       [2]uint32{6, 10},
		PotentialWeights: [10]uint32{0, 2, 10, 40, 60, 80, 60, 40, 2, 1},
		ProductId:        "player_tier_1",
		GoalKeepersCount: 3,
		DefendersCount:   6,
		MidfieldersCount: 6,
		AttackersCount:   6,
	}
	tiers = append(tiers, tier)

	// Tier2
	// - one reroll: 21 players. So if we want 1 player of 100% potential per week,
	// we need 1 in 11*2*7 = 154
	tier = WorldPlayersTier{
		LevelRange:       [2]uint32{20, 26},
		PotentialWeights: [10]uint32{0, 2, 5, 10, 20, 30, 20, 10, 2, 1},
		ProductId:        "player_tier_2",
		GoalKeepersCount: 2,
		DefendersCount:   3,
		MidfieldersCount: 3,
		AttackersCount:   3,
	}
	tiers = append(tiers, tier)

	// Tier3
	tier = WorldPlayersTier{
		LevelRange:       [2]uint32{40, 46},
		PotentialWeights: [10]uint32{0, 2, 4, 6, 15, 20, 15, 6, 2, 1},
		ProductId:        "player_tier_3",
		GoalKeepersCount: 1,
		DefendersCount:   2,
		MidfieldersCount: 2,
		AttackersCount:   2,
	}
	tiers = append(tiers, tier)

	// Tier4
	tier = WorldPlayersTier{
		LevelRange:       [2]uint32{70, 76},
		PotentialWeights: [10]uint32{0, 1, 2, 4, 10, 15, 10, 3, 2, 1},
		ProductId:        "player_tier_4",
		GoalKeepersCount: 0,
		DefendersCount:   0,
		MidfieldersCount: 0,
		AttackersCount:   0,
	}
	randomPosPlayersCount := int64(4)
	tier = addPlayerAtRandomFieldPos(tier, seed+"salt", randomPosPlayersCount)
	tiers = append(tiers, tier)

	// Tier5
	tier = WorldPlayersTier{
		LevelRange:       [2]uint32{100, 104},
		PotentialWeights: [10]uint32{0, 2, 4, 8, 16, 8, 4, 2, 1, 0},
		ProductId:        "player_tier_5",
		GoalKeepersCount: 0,
		DefendersCount:   0,
		MidfieldersCount: 0,
		AttackersCount:   0,
	}
	randomPosPlayersCount = int64(1)
	tier = addPlayerAtRandomFieldPos(tier, seed+"salt", randomPosPlayersCount)
	tiers = append(tiers, tier)

	return tiers
}
