package service

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (s *Service) process() (bool, error) {

	var err error
	var legueCount uint64

	fmt.Println("=========================================================")
	fmt.Println(s.stakers.Info())

	if legueCount, err = s.lionel.LeagueCount(); err != nil {
		return false, err
	}

	for i := uint64(0); i < legueCount; i++ {
		log.Info("Scanning league #", i)

		// find free staker
		staker, err := s.stakers.NextFreeStaker()
		if err != nil {
			return false, err
		} else if staker == nil {
			log.Info("No free stakers available to process events")
			break
		}

		// process update events ---------------------------------------
		canBeUpdated, err := s.lionel.CanLeagueBeUpdated(i)
		if err == nil {
			if canBeUpdated {
				if err := s.lionel.Update(*staker, i); err != nil {
					log.Error("Failed Update: ", err)
				}
			}
		} else {
			log.Error("Failed CanLeagueBeUpdated: ", err)
		}

		// process challange events ---------------------------------------
		canLeagueBeChallanged, err := s.lionel.CanLeagueBeChallanged(i)
		if err == nil {
			if canLeagueBeChallanged {
				if err := s.lionel.Challange(*staker, i); err != nil {
					log.Error("Failed Challange: ", err)
				}
			}
		} else {
			log.Error("Failed CanLeagueBeChallanged: ", err)
		}
	}

	//  -------------------------------------------------------------------
	for _, staker := range s.stakers.Members() {
		needsTouch, err := s.stakers.NeedsTouch(staker.Address)
		if err != nil {
			log.Error("Failed NeedsTouch: ", err)
		} else if needsTouch {
			if err = s.stakers.Touch(staker.Address); err != nil {
				log.Error("Failed Touch: ", err)
			}
			continue
		}
		needsResolve, err := s.stakers.NeedsResolve(staker.Address)
		if err != nil {
			log.Error("Failed NeedsResolve: ", err)
		} else if needsResolve {
			if err = s.stakers.Resolve(staker.Address); err != nil {
				log.Error("Failed Resolve: ", err)
			}
			continue
		}
	}

	return true, nil
}
