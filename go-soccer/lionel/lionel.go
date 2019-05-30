package lionel

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	cfg "github.com/freeverseio/go-soccer/config"
	"github.com/freeverseio/go-soccer/eth"
	"github.com/freeverseio/go-soccer/stakers"
	"github.com/freeverseio/go-soccer/storage"
	sto "github.com/freeverseio/go-soccer/storage"
)

type Lionel struct {
	storage *sto.Storage
	assets  *eth.Contract
	leagues *eth.Contract
	state   *eth.Contract

	stakers *stakers.Stakers
}

func New(web3 *eth.Web3Client, storage *sto.Storage, stakers *stakers.Stakers) (*Lionel, error) {

	// load assets
	assetsAbi, err := abi.JSON(strings.NewReader(assetsAbiJson))
	if err != nil {
		return nil, err
	}
	assetsAddress := common.HexToAddress(cfg.C.Contracts.AssetsAddress)
	assets, err := eth.NewContract(web3, &assetsAbi, nil, &assetsAddress)

	// load leagues
	leaguesAbi, err := abi.JSON(strings.NewReader(leaguesAbiJson))
	if err != nil {
		return nil, err
	}
	leaguesAddress := common.HexToAddress(cfg.C.Contracts.LeaguesAddress)
	leagues, err := eth.NewContract(web3, &leaguesAbi, nil, &leaguesAddress)

	// load state
	stateAbi, err := abi.JSON(strings.NewReader(stateAbiJson))
	if err != nil {
		return nil, err
	}
	stateAddress := common.HexToAddress(cfg.C.Contracts.StateAddress)
	state, err := eth.NewContract(web3, &stateAbi, nil, &stateAddress)

	if err != nil {
		return nil, err
	}

	return &Lionel{
		stakers: stakers,
		storage: storage,

		assets:  assets,
		leagues: leagues,
		state:   state,
	}, nil
}

func (l *Lionel) Update(staker common.Address, leagueIdx uint64) error {

	var err error

	var teamIdxs []*big.Int
	if err := l.leagues.Call(&teamIdxs, "getTeams", big.NewInt(int64(leagueIdx))); err != nil {
		return err
	}
	var countLeagueDays *big.Int
	if err := l.leagues.Call(&countLeagueDays, "countLeagueDays", big.NewInt(int64(leagueIdx))); err != nil {
		return err
	}

	userActions := []sto.UserActions{}
	for teamNo := 0; teamNo < len(teamIdxs); teamNo++ {
		tactics := [][3]uint8{}
		for dayNo := 0; dayNo < int(countLeagueDays.Uint64()); dayNo++ {
			tactics = append(tactics, [3]uint8{4, 4, 2})
		}
		userActions = append(userActions, sto.UserActions{
			Tactics: tactics,
		})
	}

	isLier, err := l.stakers.IsLier(staker)
	if err != nil {
		return err
	}

	res, err := l.ComputeLeague(
		big.NewInt(int64(leagueIdx)),
		teamIdxs,
		userActions,
		isLier,
	)
	if err != nil {
		return err
	}

	stk := l.stakers.Get(staker)

	isLieHardCoded := false // TODO
	tx, _, err := l.leagues.SendTransactionSyncWithClient(
		stk.Client, nil, 0,
		"updateLeague",
		big.NewInt(int64(leagueIdx)),
		res.initStatesHash,
		res.statesAtMatchdayHashes,
		res.scores,
		isLieHardCoded,
	)

	fmt.Printf("updateLeague leagueIdx: %v\n", leagueIdx)
	fmt.Printf("updateLeague initStatesHash: %v\n", hex.EncodeToString(res.initStatesHash[:]))
	for i, v := range res.statesAtMatchdayHashes {
		fmt.Printf("updateLeague statesAtMatchdayHashes[%v]=%v\n", i, hex.EncodeToString(v[:]))
	}
	for i, v := range res.scores {
		fmt.Printf("updateLeague scores[%v]=%v\n", i, v)
	}

	if err == nil {
		log.WithField("tx", tx.Hash().Hex()).Info("  League ", leagueIdx, " : updating lier=", isLier)
	} else {
		log.Error("  League ", leagueIdx, " : update failed")
	}
	return err
}

func (l *Lionel) Challange(staker common.Address, leagueNo uint64) error {
	/*s
	stk := l.stakers.Get(staker)

	var hash common.Hash
	tx, _, err := l.contract.SendTransactionSyncWithClient(stk.Client, nil, 0,
		"challange",
		big.NewInt(int64(leagueNo)), hash)

	if err == nil {
		log.WithField("tx", tx.Hash().Hex()).Info("  League ", leagueNo, " : challanging")
	} else {
		log.Error("  League ", leagueNo, " : challange failed")
	}
	return err
	*/
	return fmt.Errorf("Unimplemented")
}

func (l *Lionel) LeaguesCount() (uint64, error) {
	var leaguesCount *big.Int
	if err := l.leagues.Call(&leaguesCount, "leaguesCount"); err != nil {
		return 0, err
	}
	return leaguesCount.Uint64(), nil
}

func (l *Lionel) CanLeagueBeUpdated(leagueNo uint64) (bool, error) {

	leagueNoNum := big.NewInt(int64(leagueNo))

	var hasFinished bool
	if err := l.leagues.Call(&hasFinished, "hasFinished", leagueNoNum); err != nil {
		return false, err
	}
	var isUpdated bool
	if err := l.leagues.Call(&isUpdated, "isUpdated", leagueNoNum); err != nil {
		return false, err
	}
	var isVerified bool
	if err := l.leagues.Call(&isVerified, "isVerified", leagueNoNum); err != nil {
		return false, err
	}
	return hasFinished && !isUpdated && !isVerified, nil
}

func (l *Lionel) CanLeagueBeChallanged(leagueNo uint64) (bool, error) {
	return false, nil
	/*
		var canLeagueBeChallanged bool
		if err := l.contract.Call(&canLeagueBeChallanged, "canLeagueBeChallanged", big.NewInt(int64(leagueNo))); err != nil {
			return false, err
		}
		return canLeagueBeChallanged, nil
	*/
}

func (l *Lionel) generateTeamState(teamId *big.Int) ([]*big.Int, error) {
	// let teamState = await state.teamStateCreate().should.be.fulfilled;
	var teamState []*big.Int
	if err := l.state.Call(&teamState, "teamStateCreate"); err != nil {
		return nil, err
	}

	// const playersIds = await assets.getTeamPlayerIds(id).should.be.fulfilled;
	var playersIds [11]*big.Int
	if err := l.assets.Call(&playersIds, "getTeamPlayerIds", teamId); err != nil {
		return nil, err
	}

	// for (let i = 0; i < playersIds.length; i++) {
	//	const playerState = await assets.getPlayerState(playersIds[i]).should.be.fulfilled;
	//	teamState = await state.teamStateAppend(teamState, playerState).should.be.fulfilled;
	// }
	for _, playerId := range playersIds {
		var playerState *big.Int

		if err := l.assets.Call(&playerState, "getPlayerState", playerId); err != nil {
			return nil, err
		}
		if err := l.state.Call(&teamState, "teamStateAppend", teamState, playerState); err != nil {
			return nil, err
		}
	}
	return teamState, nil
}

func (l *Lionel) prepareMatchdayHashes(statesAtMatchday [][]*big.Int) ([][32]byte, error) {
	// let result = [];
	// for (let i = 0; i < statesAtMatchday.length; i++) {
	//	const state = statesAtMatchday[i];
	//	const hash = await leagues.hashDayState(state).should.be.fulfilled;
	//	result.push(hash);
	//}
	//return result;
	result := [][32]byte{}
	for _, stateAtMatchday := range statesAtMatchday {
		var hash [32]byte
		if err := l.leagues.Call(&hash, "hashDayState", stateAtMatchday); err != nil {
			return nil, err
		}
		result = append(result, hash)
	}
	return result, nil
}

type LeagueResult struct {
	initStatesHash         [32]byte
	statesAtMatchdayHashes [][32]byte
	scores                 []uint16
}

func (l *Lionel) ComputeLeague(leagueIdx *big.Int, teamIdxs []*big.Int, actionsPerDay []storage.UserActions, lier bool) (*LeagueResult, error) {

	// compute leagueState at beginning of the league
	var initLeagueState []*big.Int

	if err := l.state.Call(&initLeagueState, "leagueStateCreate"); err != nil {
		return nil, err
	}

	for _, teamIdx := range teamIdxs {
		var teamState []*big.Int
		var err error
		if teamState, err = l.generateTeamState(teamIdx); err != nil {
			return nil, err
		}
		if err := l.state.Call(&initLeagueState, "leagueStateAppend", initLeagueState, teamState); err != nil {
			return nil, err
		}
	}

	leagueState := make([]*big.Int, len(initLeagueState))
	copy(leagueState, initLeagueState)

	scores := []uint16{}
	if err := l.leagues.Call(&scores, "scoresCreate"); err != nil {
		return nil, err
	}

	statesAtMatchday := [][]*big.Int{}

	for day, dayUserActions := range actionsPerDay {
		// update league state given day user actions
		type ComputeDateResultType struct {
			Scores           []uint16
			FinalLeagueState []*big.Int
		}
		var computeDateResult ComputeDateResultType
		if err := l.leagues.Call(&computeDateResult, "computeDay", leagueIdx, big.NewInt(int64(day)), leagueState, dayUserActions.Tactics); err != nil {
			return nil, err
		}

		leagueState = computeDateResult.FinalLeagueState

		// apend scores
		statesAtMatchday = append(statesAtMatchday, computeDateResult.FinalLeagueState)

		if err := l.leagues.Call(&scores, "scoresConcat", scores, computeDateResult.Scores); err != nil {
			return nil, err
		}
	}

	// let updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
	// updated.should.be.equal(false);

	var isUpdated bool
	if err := l.leagues.Call(&isUpdated, "isUpdated", leagueIdx); err != nil {
		return nil, err
	}
	if isUpdated {
		return nil, fmt.Errorf("leagues.isUpdated(leagueIdx) failed")
	}

	// const initStatesHash = await leagues.hashInitState(initPlayerStates).should.be.fulfilled;
	// const statesAtMatchdayHashes = await prepareMatchdayHashes(statesAtMatchday);
	var initStatesHash [32]byte
	if err := l.leagues.Call(&initStatesHash, "hashInitState", initLeagueState); err != nil {
		return nil, err
	}

	if lier {
		statesAtMatchday[0][0] = statesAtMatchday[0][0].Add(statesAtMatchday[0][0], big.NewInt(1))
	}

	statesAtMatchdayHashes, err := l.prepareMatchdayHashes(statesAtMatchday)
	if err != nil {
		return nil, err
	}

	lr := LeagueResult{
		initStatesHash:         initStatesHash,
		statesAtMatchdayHashes: statesAtMatchdayHashes,
		scores:                 scores,
	}

	return &lr, nil
}
