package service

import (
	log "github.com/sirupsen/logrus"
)

func (s *Service) processLionel() (bool, error) {

	var err error
	var legueCount uint64

	if legueCount, err = s.lionel.LeagueCount(); err != nil {
		return false, err
	}
	log.Info("Scanning ", legueCount, " leagues...")

	for i := uint64(0); i < legueCount; i++ {

		canBeUpdated, err := s.lionel.CanLeagueBeUpdated(i)
		if err != nil {
			return false, err
		}

		if canBeUpdated {
			if err := s.lionel.Update(i); err != nil {
				return false, err
			}
		}

		canLeagueBeChallanged, err := s.lionel.CanLeagueBeChallanged(i)
		if err != nil {
			return false, err
		}
		if canLeagueBeChallanged {
			if err := s.lionel.Challange(i); err != nil {
				return false, err
			}
		}

	}

	return true, nil
}
