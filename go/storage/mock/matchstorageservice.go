package mock

import (
	"github.com/freeverseio/crypto-soccer/go/storage"
)

type MatchStorageService struct {
	MatchesByTimezoneFunc func(timezone uint8) ([]storage.Match, error)
}

func (b MatchStorageService) MatchesByTimezone(timezone uint8) ([]storage.Match, error) {
	return b.MatchesByTimezoneFunc(timezone)
}
