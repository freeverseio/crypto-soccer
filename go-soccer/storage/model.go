package storage

import "math/big"

type GlobalsEntry struct {
}

type LeagueEntry struct {
}

type UserActionsEntry struct {
	eagueIdx      uint64
	TeamIdxs      []*big.Int
	ActionsPerDay []UserActions
}

type UserActions struct {
	Tactics [][3]uint8
}
