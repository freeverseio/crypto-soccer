package worldplayer

import (
	"math/big"
	"strconv"
)

type WorldPlayersTier struct {
	Value            int64
	MaxPotential     uint8
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
	// - maxPotential 80%
	tier := WorldPlayersTier{
		Value:            1200,
		MaxPotential:     7,
		ProductId:        "player_tier_1",
		GoalKeepersCount: 3,
		DefendersCount:   6,
		MidfieldersCount: 6,
		AttackersCount:   6,
	}
	tiers = append(tiers, tier)

	// Tier2:
	// - a fixed number of players, with field position distributed randomly
	// - maxPotential 90%
	tier = WorldPlayersTier{
		Value:            2200,
		MaxPotential:     8,
		ProductId:        "player_tier_2",
		GoalKeepersCount: 0,
		DefendersCount:   0,
		MidfieldersCount: 0,
		AttackersCount:   0,
	}
	randomPosPlayersCount := int64(4)
	tier = addPlayerAtRandomFieldPos(tier, seed, randomPosPlayersCount)
	tiers = append(tiers, tier)

	// Tier3
	// - a fixed number of players, with field position distributed randomly
	// - field position distributed randomly
	// - maxPotential 100%
	tier = WorldPlayersTier{
		Value:            3500,
		MaxPotential:     9,
		ProductId:        "player_tier_3",
		GoalKeepersCount: 0,
		DefendersCount:   0,
		MidfieldersCount: 0,
		AttackersCount:   0,
	}
	randomPosPlayersCount = int64(3)
	tier = addPlayerAtRandomFieldPos(tier, seed+"salt", randomPosPlayersCount)
	tiers = append(tiers, tier)

	return tiers
}
