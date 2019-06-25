package memory

import "errors"

type Team struct {
	Name string
}

type MemoryStorage struct {
	teams map[uint64]Team
}

func New() *MemoryStorage {
	return &MemoryStorage{
		teams: make(map[uint64]Team),
	}
}

func (m *MemoryStorage) TeamAdd(ID uint64, name string) error {
	m.teams[ID] = Team{name}
	return nil
}

func (m *MemoryStorage) TeamCount() (uint64, error) {
	return uint64(len(m.teams)), nil
}

func (m *MemoryStorage) GetTeam(id uint64) (Team, error) {
	team := m.teams[id]
	if team.Name == "" {
		return team, errors.New("unexistent team")
	}
	return team, nil
}
