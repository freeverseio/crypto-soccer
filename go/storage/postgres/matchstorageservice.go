package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/storage"

	log "github.com/sirupsen/logrus"
)

type MatchStorageService struct {
}

func NewMatchStorageService(tx *sql.Tx) *MatchStorageService {
	return &MatchStorageService{}
}

// func (b MatchStorageService) Insert(match storage.Match) error {
// 	if b.matches[match.TimezoneIdx] == nil {
// 		b.matches[match.TimezoneIdx] = make(map[uint32]map[uint32]storage.Match)
// 	}
// 	if b.matches[match.TimezoneIdx][match.CountryIdx] == nil {
// 		b.matches[match.TimezoneIdx][match.CountryIdx] = make(map[uint32]storage.Match)
// 	}
// 	b.matches[match.TimezoneIdx][match.CountryIdx][match.LeagueIdx] = match
// 	return nil
// }

func (b MatchStorageService) MatchesByTimezone(timezone uint8) ([]storage.Match, error) {
	log.Warning("MatchesByTimezone not implemented")
	return []storage.Match{}, nil
}
