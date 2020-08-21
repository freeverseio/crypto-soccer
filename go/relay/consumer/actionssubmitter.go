package consumer

import (
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/useractions/orgmap"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/useractions"
	"github.com/freeverseio/crypto-soccer/go/useractions/postgres"
)

type ActionsSubmitter struct {
	client                    *ethclient.Client
	contracts                 contracts.Contracts
	auth                      *bind.TransactOpts
	useractionsPublishService useractions.UserActionsPublishService
}

// *****************************************************************************
// public
// *****************************************************************************

func NewActionsSubmitter(
	client *ethclient.Client,
	auth *bind.TransactOpts,
	contracts contracts.Contracts,
	useractionsPublishService useractions.UserActionsPublishService,
) *ActionsSubmitter {
	return &ActionsSubmitter{
		client,
		contracts,
		auth,
		useractionsPublishService,
	}
}

func (p *ActionsSubmitter) Process(tx *sql.Tx) error {
	nextUpdate, err := p.NextUpdateSinceEpochSec()
	now := NowSinceEpochSec()
	if now < nextUpdate {
		log.Infof("Countdown: %v (%v - now: %v)", nextUpdate-now, nextUpdate, now)
		return nil
	}
	currentVerse, err := p.contracts.Updates.GetCurrentVerse(&bind.CallOpts{})
	if err != nil {
		return err
	}
	nextToUpdate, err := p.contracts.Updates.NextTimeZoneToUpdate(&bind.CallOpts{})
	if err != nil {
		return err
	}
	log.Infof("Staring process of verse %v, timezone %v, day %v, turn %v", currentVerse, nextToUpdate.Tz, nextToUpdate.Day, nextToUpdate.TurnInDay)
	upcomingUserActions := &useractions.UserActions{}
	if nextToUpdate.Tz == 0 {
		log.Info("Timezone 0 ... skipping user actions")
	} else if nextToUpdate.TurnInDay <= 1 {
		useractionsStorageService := postgres.NewUserActionsStorageService(tx)
		if upcomingUserActions, err = useractionsStorageService.UserActionsByTimezone(int(nextToUpdate.Tz)); err != nil {
			return err
		}
	}
	root, err := upcomingUserActions.Root()
	if err != nil {
		return err
	}
	var cid string
	if nextToUpdate.Day == 7 && nextToUpdate.TurnInDay == 1 {
		//TODO: generate orgmapdenylist or read it from somewhere, maybe db??
		upcomingUserActions.OrgMapDenyList = make([]orgmap.OrgMapDenyList, 0)
		cid, err = p.useractionsPublishService.Publish(*upcomingUserActions)
	} else {
		cid, err = p.useractionsPublishService.Publish(*upcomingUserActions)
	}
	if err != nil {
		return err
	}
	log.Infof("[relay] submitActionsRoot root: 0x%v, cid: %v", hex.EncodeToString(root[:]), cid)
	transaction, err := p.contracts.Updates.SubmitActionsRoot(p.auth, root, root, root, 3, cid)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(p.client, transaction, 30) // TODO make timeout configurable
	if err != nil {
		return err
	}
	return nil
}

func (p *ActionsSubmitter) NextUpdateSinceEpochSec() (int64, error) {
	secs, err := p.contracts.Updates.GetNextVerseTimestamp(nil)
	if err != nil {
		return 0, err
	}
	return secs.Int64(), nil
}

func NowSinceEpochSec() int64 {
	now := time.Now()
	secs := now.Unix()
	return secs
}
