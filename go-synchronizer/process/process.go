package process

import (
	"context"
	"errors"

	//"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"

	//"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/market"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/utils"
	log "github.com/sirupsen/logrus"
)

type EventProcessor struct {
	usesGanache bool
	client      *ethclient.Client
	db          *storage.Storage
	market      *market.Market
	assets      *assets.Assets
}

// *****************************************************************************
// public
// *****************************************************************************

// NewEventProcessor creates a new struct for scanning and storing crypto soccer events
func NewEventProcessor(client *ethclient.Client, db *storage.Storage, market *market.Market, assets *assets.Assets) *EventProcessor {
	return &EventProcessor{false, client, db, market, assets}
}

// NewGanacheEventProcessor creates a new struct for scanning and storing crypto soccer events from a ganache client
func NewGanacheEventProcessor(client *ethclient.Client, db *storage.Storage, market *market.Market, assets *assets.Assets) *EventProcessor {
	return &EventProcessor{true, client, db, market, assets}
}

// Process processes all scanned events and stores them into the database db
func (p *EventProcessor) Process() error {
	opts := p.nextRange()

	if opts == nil {
		log.Info("No new blocks to scan.")
		return nil
	}

	log.WithFields(log.Fields{
		"start": opts.Start,
		"end":   *opts.End,
	}).Info("Syncing ...")

	if events, err := p.scanDivisionCreation(opts); err != nil {
		return err
	} else {
		err = p.storeDivisionCreation(events)
		if err != nil {
			return err
		}
	}

	if events, err := p.scanTeamTransfer(opts); err != nil {
		return err
	} else {
		for _, event := range events { // TODO: next part to be recoded
			// _, blockNumber, err := p.getTimeOfEvent(event.Raw)
			// if err != nil {
			// 	return err
			// }
			teamID := event.TeamId
			newOwner := event.To.String()
			team, err := p.db.GetTeam(teamID)
			if err != nil {
				return err
			}
			// team.State.BlockNumber = blockNumber
			team.State.Owner = newOwner
			err = p.db.TeamUpdate(teamID, team.State)
			if err != nil {
				return err
			}
		}
	}

	//if events, err := p.scanPlayerTransfer(opts); err != nil {
	//	return err
	//} else {
	//	for _, event := range events { // TODO: next part to be recoded
	//		_, blockNumber, err := p.getTimeOfEvent(event.Raw)
	//		if err != nil {
	//			return err
	//		}
	//		playerId := event.PlayerId.Uint64()
	//		toTeamId := event.ToTeamId.Uint64()
	//		player, err := p.db.GetPlayer(playerId)
	//		if err != nil {
	//			return err
	//		}
	//		player.State.BlockNumber = blockNumber
	//		player.State.TeamId = toTeamId
	//		err = p.db.PlayerStateUpdate(playerId, player.State)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}

	//if p.leagues != nil {
	//	if events, err := p.scanLeagueCreated(opts); err != nil {
	//		return err
	//	} else {
	//		for _, event := range events { // TODO: next part to be recoded
	//			p.db.LeagueAdd(storage.League{
	//				Id: event.LeagueId.Uint64(),
	//			})
	//			// log.Info(
	//			// 	"Found league ", event.LeagueId.Int64(),
	//			// 	"\n\tdays: ", p.getLeagueDaysCount(event.LeagueId),
	//			// 	"\n\tfinished: ", p.hasLeagueFinished(event.LeagueId),
	//			// 	"\n\tupdated: ", p.isLeagueUpdated(event.LeagueId),
	//			// )
	//		}
	//	}
	//} else {
	//	log.Warn("Contract leagues not set. Skipping scanning for leagues")
	//}

	// store the last block that was scanned
	p.db.SetBlockNumber(*opts.End)
	return nil
}

// *****************************************************************************
// private
// *****************************************************************************
func (p *EventProcessor) nextRange() *bind.FilterOpts {
	start := p.dbLastBlockNumber()
	if start != 0 {
		// unless this is the very first execution,
		// the block number that is stored in the db
		// was already scanned. We are interested in
		// the next block
		if start < math.MaxUint64 {
			start += 1
		} else {
			log.Error("Block range overflow")
			return nil
		}
	}
	end := p.clientLastBlockNumber()
	if start > end {
		return nil
	}
	return &bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	}
}

func (p *EventProcessor) clientLastBlockNumber() uint64 {
	if p.client == nil {
		log.Warn("Client is nil. Returning 0 as last block.")
		return 0
	}
	header, err := p.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Warn("Could not get blockchain last block")
		return 0
	}
	return header.Number.Uint64()
}
func (p *EventProcessor) dbLastBlockNumber() uint64 {
	storedLastBlockNumber, err := p.db.GetBlockNumber()
	if err != nil {
		log.Warn("Could not get database last block")
		return 0
	}
	return storedLastBlockNumber
}
func (p *EventProcessor) getTimeOfEvent(eventRaw types.Log) (uint64, uint64, error) {
	if p.usesGanache {
		return eventRaw.BlockNumber, eventRaw.BlockNumber, nil
	}
	block, err := p.client.BlockByHash(context.Background(), eventRaw.BlockHash)
	if err != nil {
		return 0, 0, err
	}
	return block.Time(), eventRaw.BlockNumber, nil
}

func (p *EventProcessor) scanDivisionCreation(opts *bind.FilterOpts) ([]assets.AssetsDivisionCreation, error) {
	if opts == nil {
		opts = &bind.FilterOpts{Start: 0}
	}
	iter, err := p.assets.FilterDivisionCreation(opts)
	if err != nil {
		return nil, err
	}

	events := []assets.AssetsDivisionCreation{}

	for iter.Next() {
		events = append(events, *(iter.Event))
	}
	return events, nil
}

func (p *EventProcessor) storeDivisionCreation(events []assets.AssetsDivisionCreation) error {
	for _, event := range events {
		log.Info(
			"\ntime zone: ", event.Timezone,
			"\nCountry idx: ", event.CountryIdxInTZ.Uint64(),
			"\ndivision idx: ", event.DivisionIdxInCountry.Uint64())
		if event.CountryIdxInTZ.Uint64() == 0 {
			if err := p.db.TimezoneCreate(storage.Timezone{event.Timezone}); err != nil {
				return err
			}
		}
		if event.DivisionIdxInCountry.Uint64() == 0 {
			countryIdx := event.CountryIdxInTZ.Uint64()
			if countryIdx > 65535 {
				return errors.New("Cannot cast country idx to uint16: value too large")
			}
			if err := p.db.CountryCreate(storage.Country{event.Timezone, uint16(countryIdx)}); err != nil {
				return err
			}
			if err := p.storeTeamsForNewDivision(event.Timezone, event.CountryIdxInTZ, event.DivisionIdxInCountry); err != nil {
				return err
			}
		}
	}
	return nil
}
func (p *EventProcessor) storeTeamsForNewDivision(timezone uint8, countryIdx *big.Int, divisionIdxInCountry *big.Int) error {
	opts := &bind.CallOpts{}

	LEAGUES_PER_DIV, err := p.assets.LEAGUESPERDIV(opts)
	if err != nil {
		return err
	}
	leagueIdxBegin := divisionIdxInCountry.Int64() * int64(LEAGUES_PER_DIV)
	leagueIdxEnd := leagueIdxBegin + int64(LEAGUES_PER_DIV)

	TEAMS_PER_LEAGUE, err := p.assets.TEAMSPERLEAGUE(opts)
	if err != nil {
		return err
	}

	for leagueIdx := leagueIdxBegin; leagueIdx < leagueIdxEnd; leagueIdx++ {
		teamIdxBegin := leagueIdx
		teamIdxEnd := teamIdxBegin + int64(TEAMS_PER_LEAGUE)
		for teamIdx := teamIdxBegin; teamIdx < teamIdxEnd; teamIdx++ {
			if teamId, e := p.assets.EncodeTZCountryAndVal(opts, timezone, countryIdx, big.NewInt(teamIdx)); e != nil {
				return e
			} else {
				if teamOwner, e := p.assets.GetOwnerTeam(opts, teamId); e != nil {
					return e
				} else if e := p.db.TeamCreate(
					storage.Team{
						teamId,
						timezone,
						uint16(countryIdx.Uint64()),
						storage.TeamState{teamOwner.Hex(), uint8(leagueIdx)}},
				); e != nil {
					return e
				} else if e := p.storeVirtualPlayersForTeam(opts, teamId, timezone, countryIdx, teamIdx); e != nil {
					return e
				}
			}
		}
	}
	return err
}
func (p *EventProcessor) storeVirtualPlayersForTeam(opts *bind.CallOpts, teamId *big.Int, timezone uint8, countryIdx *big.Int, teamIdxInCountry int64) error {
	PLAYERS_PER_TEAM_INIT, err := p.assets.PLAYERSPERTEAMINIT(opts)
	if err != nil {
		return err
	}
	begin := teamIdxInCountry * int64(PLAYERS_PER_TEAM_INIT)
	end := begin + int64(PLAYERS_PER_TEAM_INIT)

	SK_SHO := uint8(0)
	SK_SPE := uint8(0)
	SK_PAS := uint8(0)
	SK_DEF := uint8(0)
	SK_END := uint8(0)

	SK_SHO, err = p.assets.SKSHO(opts)
	if err != nil {
		return err
	}
	SK_SPE, err = p.assets.SKSPE(opts)
	if err != nil {
		return err
	}
	SK_PAS, err = p.assets.SKPAS(opts)
	if err != nil {
		return err
	}
	SK_DEF, err = p.assets.SKDEF(opts)
	if err != nil {
		return err
	}
	SK_END, err = p.assets.SKEND(opts)
	if err != nil {
		return err
	}

	for i := begin; i < end; i++ {
		if playerId, e := p.assets.EncodeTZCountryAndVal(opts, timezone, countryIdx, big.NewInt(i)); e != nil {
			return e
		} else if skills, e := p.getPlayerSkillsAtBirth(opts, playerId); e != nil {
			return e
		} else if preferredPosition, e := p.getPlayerPreferredPosition(opts, playerId); e != nil {
			return e
		} else if e := p.db.PlayerCreate(
			storage.Player{
				playerId,
				preferredPosition,
				storage.PlayerState{ // TODO: storage should use same skill ordering as BC
					TeamId:    teamId,
					Defence:   uint64(skills[SK_DEF]), // TODO: type should be uint16
					Speed:     uint64(skills[SK_SPE]),
					Pass:      uint64(skills[SK_PAS]),
					Shoot:     uint64(skills[SK_SHO]),
					Endurance: uint64(skills[SK_END]),
				}},
		); e != nil {
			return e
		}
	}
	return err
}

func (p *EventProcessor) getPlayerSkillsAtBirth(opts *bind.CallOpts, playerId *big.Int) ([5]uint16, error) {
	if skills, err := p.assets.GetPlayerSkillsAtBirth(opts, playerId); err != nil {
		return [5]uint16{0, 0, 0, 0, 0}, err
	} else {
		return p.assets.GetSkillsVec(opts, skills)
	}
}

func (p *EventProcessor) getPlayerPreferredPosition(opts *bind.CallOpts, playerId *big.Int) (string, error) {
	if encodedSkills, err := p.assets.GetPlayerSkillsAtBirth(opts, playerId); err != nil {
		return "", err
	} else if forwardness, err := p.assets.GetForwardness(opts, encodedSkills); err != nil {
		return "", err
	} else if leftishness, err := p.assets.GetLeftishness(opts, encodedSkills); err != nil {
		return "", err
	} else {
		if forwardness.Uint64() > 255 {
			return "", errors.New("Cannot cast forwardness to uint8: value too large")
		} else if leftishness.Uint64() > 255 {
			return "", errors.New("Cannot cast leftishness to uint8: value too large")
		}
		return utils.PreferredPosition(uint8(forwardness.Uint64()), uint8(leftishness.Uint64()))
	}
}

//func (p *EventProcessor) storeTeamCreated(events []assets.AssetsTeamCreated) error {
//	for _, event := range events {
//		if name, err := p.assets.GetTeamName(nil, event.Id); err != nil {
//			return err
//		} else if owner, err := p.assets.GetTeamOwner(nil, event.Id); err != nil {
//			return err
//		} else if blockTime, blockNumber, err := p.getTimeOfEvent(event.Raw); err != nil {
//			return err
//		} else if err := p.db.TeamAdd(storage.Team{
//			Id:                event.Id.Uint64(),
//			Name:              name,
//			CreationTimestamp: blockTime,
//			CountryId:         1, // TODO: get it from blockchain
//			State: storage.TeamState{
//				BlockNumber:          blockNumber,
//				Owner:                owner.Hex(),
//				CurrentLeagueId:      1, // TODO: uint64
//				PosInCurrentLeagueId: 0, // TODO: uint64
//				PrevLeagueId:         0, // TODO: uint64
//				PosInPrevLeagueId:    0, // TODO: uint64
//			},
//		}); err != nil {
//			return err
//		}
//		if err := p.storeVirtualPlayers(event.Id); err != nil {
//			return err
//		}
//	}
//	return nil
//}

//func (p *EventProcessor) scanTeamCreated(opts *bind.FilterOpts) ([]assets.AssetsTeamCreated, error) {
//	if opts == nil {
//		opts = &bind.FilterOpts{Start: 0}
//	}
//	iter, err := p.assets.FilterTeamCreated(opts)
//	if err != nil {
//		return nil, err
//	}
//
//	events := []assets.AssetsTeamCreated{}
//
//	for iter.Next() {
//		events = append(events, *(iter.Event))
//	}
//	return events, nil
//}
func (p *EventProcessor) scanTeamTransfer(opts *bind.FilterOpts) ([]assets.AssetsTeamTransfer, error) {
	if opts == nil {
		opts = &bind.FilterOpts{Start: 0}
	}
	iter, err := p.assets.FilterTeamTransfer(opts)
	if err != nil {
		return nil, err
	}

	events := []assets.AssetsTeamTransfer{}

	for iter.Next() {
		events = append(events, *(iter.Event))
	}
	return events, nil
}

//func (p *EventProcessor) scanPlayerTransfer(opts *bind.FilterOpts) ([]assets.AssetsPlayerTransfer, error) {
//	if opts == nil {
//		opts = &bind.FilterOpts{Start: 0}
//	}
//	iter, err := p.assets.FilterPlayerTransfer(opts)
//	if err != nil {
//		return nil, err
//	}
//
//	events := []assets.AssetsPlayerTransfer{}
//
//	for iter.Next() {
//		events = append(events, *(iter.Event))
//	}
//	return events, nil
//}

//func (p *EventProcessor) storeVirtualPlayers(teamId *big.Int) error {
//	// TODO: move to a single run place ...  constructor
//	nPlayersAtCreation, err := p.assets.PLAYERSPERTEAMINIT(&bind.CallOpts{})
//	if err != nil {
//		return err
//	}
//
//	for i := 0; i < int(nPlayersAtCreation); i++ {
//		if id, err := p.assets.GenerateVirtualPlayerId(&bind.CallOpts{}, teamId, uint8(i)); err != nil {
//			return err
//		} else if state, err := p.assets.GenerateVirtualPlayerState(&bind.CallOpts{}, id); err != nil {
//			return err
//		} else {
//			if skills, err := p.states.GetSkillsVec(&bind.CallOpts{}, state); err != nil {
//				return err
//			} else {
//				player := storage.Player{
//					Id:                     id.Uint64(),
//					MonthOfBirthInUnixTime: "0", // TODO
//					State: storage.PlayerState{
//						TeamId:    teamId.Uint64(),
//						State:     state.String(),
//						Defence:   uint64(skills[0]),
//						Speed:     uint64(skills[1]),
//						Pass:      uint64(skills[2]),
//						Shoot:     uint64(skills[3]),
//						Endurance: uint64(skills[4]),
//					},
//				}
//				p.db.PlayerAdd(player)
//			}
//		}
//	}
//	return nil
//}

//func (p *EventProcessor) scanLeagueCreated(opts *bind.FilterOpts) ([]leagues.LeaguesLeagueCreated, error) {
//	if opts == nil {
//		opts = &bind.FilterOpts{Start: 0}
//	}
//	iter, err := p.leagues.FilterLeagueCreated(opts)
//	if err != nil {
//		return nil, err
//	}
//
//	events := []leagues.LeaguesLeagueCreated{}
//
//	for iter.Next() {
//		events = append(events, *(iter.Event))
//	}
//	return events, nil
//}
//func (p *EventProcessor) isLeagueUpdated(leagueId *big.Int) bool {
//	result, err := p.leagues.IsUpdated(nil, leagueId)
//	if err != nil {
//		log.Fatal(err)
//		return false
//	}
//	return result
//}
//func (p *EventProcessor) hasLeagueFinished(leagueId *big.Int) bool {
//	result, err := p.leagues.HasFinished(nil, leagueId)
//	if err != nil {
//		log.Fatal(err)
//		return false
//	}
//	return result
//}
//func (p *EventProcessor) getLeagueDaysCount(leagueId *big.Int) int64 {
//	result, err := p.leagues.CountLeagueDays(nil, leagueId)
//	if err != nil {
//		log.Fatal(err)
//		return 0
//	}
//	return result.Int64()
//}
