package memory

import (
	"errors"
	"math/big"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

type MemoryStorage struct {
	teams       map[uint64]storage.Team
	blockNumber *big.Int
}

func New() *MemoryStorage {
	return &MemoryStorage{
		teams:       make(map[uint64]storage.Team),
		blockNumber: nil,
	}
}

func (m *MemoryStorage) TeamAdd(id uint64, name string) error {
	m.teams[id] = storage.Team{id, name}
	return nil
}

func (m *MemoryStorage) TeamCount() (uint64, error) {
	return uint64(len(m.teams)), nil
}

func (m *MemoryStorage) GetTeam(id uint64) (storage.Team, error) {
	team := m.teams[id]
	if team.Name == "" {
		return team, errors.New("unexistent team")
	}
	return team, nil
}

func (m *MemoryStorage) GetBlockNumber() (*big.Int, error) {
	return m.blockNumber, nil
}

func (m *MemoryStorage) SetBlockNumber(value *big.Int) error {
	m.blockNumber = value
	return nil
}
