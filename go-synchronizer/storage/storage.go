package storage

import "math/big"

type Team struct {
	Id   uint64
	Name string
}

type Storage interface {
	TeamAdd(ID uint64, name string) error
	TeamCount() (uint64, error)
	GetTeam(id uint64) (Team, error)
	GetBlockNumber() (*big.Int, error)
	SetBlockNumber(value *big.Int) error
}
