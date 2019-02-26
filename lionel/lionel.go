package lionel

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	cfg "github.com/freeverseio/go-soccer/config"
	"github.com/freeverseio/go-soccer/eth"
	sto "github.com/freeverseio/go-soccer/storage"
	log "github.com/sirupsen/logrus"
)

type Lionel struct {
	web3     *eth.Web3Client
	storage  *sto.Storage
	contract *eth.Contract
}

func New(web3 *eth.Web3Client, storage *sto.Storage) (*Lionel, error) {

	// load lionel

	lionelAbi, err := abi.JSON(strings.NewReader(lionelAbiJson))
	if err != nil {
		return nil, err
	}
	lionelAddress := common.HexToAddress(cfg.C.Contracts.LionelAddress)

	contract, err := eth.NewContract(web3, &lionelAbi, nil, &lionelAddress)
	if err != nil {
		return nil, err
	}

	return &Lionel{
		web3:     web3,
		storage:  storage,
		contract: contract,
	}, nil
}

func (l *Lionel) Update(leagueNo uint64) error {

	var hash common.Hash
	tx, _, err := l.contract.SendTransactionSync(nil, 0,
		"update",
		big.NewInt(int64(leagueNo)), &hash)

	if err == nil {
		log.WithField("tx", tx.Hash().Hex()).Info("  League ", leagueNo, " : updating")
	} else {
		log.Error("  League ", leagueNo, " : update failed")
	}
	return err
}

func (l *Lionel) Challange(leagueNo uint64) error {

	var hash common.Hash
	tx, _, err := l.contract.SendTransactionSync(nil, 0,
		"challange",
		big.NewInt(int64(leagueNo)), hash)

	if err == nil {
		log.WithField("tx", tx.Hash().Hex()).Info("  League ", leagueNo, " : challanging")
	} else {
		log.Error("  League ", leagueNo, " : challange failed")
	}
	return err
}

func (l *Lionel) LeagueCount() (uint64, error) {
	var legueCount *big.Int
	if err := l.contract.Call(&legueCount, "legueCount"); err != nil {
		return 0, err
	}
	return legueCount.Uint64(), nil
}

func (l *Lionel) CanLeagueBeUpdated(leagueNo uint64) (bool, error) {
	var canLeagueBeUpdated bool
	if err := l.contract.Call(&canLeagueBeUpdated, "canLeagueBeUpdated", big.NewInt(int64(leagueNo))); err != nil {
		return false, err
	}
	return canLeagueBeUpdated, nil
}

func (l *Lionel) CanLeagueBeChallanged(leagueNo uint64) (bool, error) {
	var canLeagueBeChallanged bool
	if err := l.contract.Call(&canLeagueBeChallanged, "canLeagueBeChallanged", big.NewInt(int64(leagueNo))); err != nil {
		return false, err
	}
	return canLeagueBeChallanged, nil
}
