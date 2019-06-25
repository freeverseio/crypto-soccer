package memory

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

type MemoryStorage struct {
	teams map[uint64]storage.Team
}

func New() *MemoryStorage {
	return &MemoryStorage{
		teams: make(map[uint64]storage.Team),
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
