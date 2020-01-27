package engine

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

type Matches []Match

func FromStorage(
	tx *sql.Tx,
	timezoneIdx uint8,
	day uint8,
) (Matches, error) {
	ms, err := storage.MatchesByTimezoneIdxAndMatchDay(tx, timezoneIdx, day)
	if err != nil {
		return nil, err
	}

	var matches Matches
	for i := range ms {
		m := ms[i]
		homeTeamPlayers, err := storage.PlayersByTeamId(tx, m.HomeTeamID)
		if err != nil {
			return nil, err
		}
		visitorTeamPlayers, err := storage.PlayersByTeamId(tx, m.VisitorTeamID)
		if err != nil {
			return nil, err
		}
		match := NewMatchFromStorage(m, homeTeamPlayers, visitorTeamPlayers)
		matches = append(matches, *match)
	}
	return matches, nil
}

func (b Matches) ToStorage(contracts *contracts.Contracts, tx *sql.Tx) error {
	for _, match := range b {
		for _, player := range match.HomeTeam.Players {
			var sPlayer storage.Player
			defence, speed, pass, shoot, endurance, _, _, err := contracts.DecodeSkills(player.Skills())
			if err != nil {
				return err
			}
			sPlayer.Defence = defence.Uint64()
			sPlayer.Speed = speed.Uint64()
			sPlayer.Pass = pass.Uint64()
			sPlayer.Shoot = shoot.Uint64()
			sPlayer.Defence = endurance.Uint64()
			sPlayer.EncodedSkills = player.Skills()
			if err = sPlayer.Update(tx); err != nil {
				return err
			}
		}
	}
	return nil
}
