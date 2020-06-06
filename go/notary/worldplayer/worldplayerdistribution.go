package worldplayer

type WorldPlayersTier struct {
	Value            int64
	MaxPotential     uint8
	ProductId        string
	GoalKeepersCount int
	DefendersCount   int
	MidfieldersCount int
	AttackersCount   int
}

func GenerateBatchDistribution() []WorldPlayersTier {
	return []WorldPlayersTier{
		WorldPlayersTier{
			Value:            1000,
			MaxPotential:     9,
			ProductId:        "player_tier_0",
			GoalKeepersCount: 1,
			DefendersCount:   2,
			MidfieldersCount: 3,
			AttackersCount:   1,
		},
		WorldPlayersTier{
			Value:            1500,
			MaxPotential:     9,
			ProductId:        "player_tier_1",
			GoalKeepersCount: 1,
			DefendersCount:   2,
			MidfieldersCount: 3,
			AttackersCount:   1,
		},
		WorldPlayersTier{
			Value:            2000,
			MaxPotential:     9,
			ProductId:        "player_tier_2",
			GoalKeepersCount: 1,
			DefendersCount:   2,
			MidfieldersCount: 3,
			AttackersCount:   1,
		}, WorldPlayersTier{
			Value:            3000,
			MaxPotential:     9,
			ProductId:        "player_tier_3",
			GoalKeepersCount: 1,
			DefendersCount:   2,
			MidfieldersCount: 3,
			AttackersCount:   1,
		}, WorldPlayersTier{
			Value:            5000,
			MaxPotential:     9,
			ProductId:        "player_tier_4",
			GoalKeepersCount: 1,
			DefendersCount:   2,
			MidfieldersCount: 3,
			AttackersCount:   1,
		}, WorldPlayersTier{
			Value:            8000,
			MaxPotential:     9,
			ProductId:        "player_tier_5",
			GoalKeepersCount: 1,
			DefendersCount:   2,
			MidfieldersCount: 3,
			AttackersCount:   1,
		}, WorldPlayersTier{
			Value:            10000,
			MaxPotential:     9,
			ProductId:        "player_tier_6",
			GoalKeepersCount: 1,
			DefendersCount:   2,
			MidfieldersCount: 3,
			AttackersCount:   1,
		},
	}
}
