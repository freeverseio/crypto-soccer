package memory

import (
	"github.com/freeverseio/crypto-soccer/go/storage"
)

type MatchStorageService struct {
	matches map[uint8]map[uint32]map[uint32]storage.Match
}

func NewMatchStorageService() *MatchStorageService {
	return &MatchStorageService{
		matches: make(map[uint8]map[uint32]map[uint32]storage.Match),
	}
}

func (b MatchStorageService) Insert(match storage.Match) error {
	b.matches[match.TimezoneIdx][match.CountryIdx][match.LeagueIdx] = match
	return nil
}

func (b MatchStorageService) MatchesByTimezone(timezone uint8) ([]storage.Match, error) {
	matches := []storage.Match{}
	for _, v := range b.matches[timezone] {
		for _, v := range v {
			matches = append(matches, v)
		}
	}
	return matches, nil
}
