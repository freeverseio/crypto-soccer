package process

import (
	"database/sql"
	"encoding/hex"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/contracts/autogenerated/updates"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/staker"
	"github.com/freeverseio/crypto-soccer/go/universe"
	"github.com/freeverseio/crypto-soccer/go/useractions"

	log "github.com/sirupsen/logrus"
)

func ConsumeActionsSubmission(
	tx *sql.Tx,
	contracts *contracts.Contracts,
	useractionsPublisher useractions.UserActionsPublishService,
	staker *staker.Staker,
	v updates.UpdatesActionsSubmission,
) error {
	log.Infof("[processor|consume] ActionsSubmission verse: %v, tz: %v, Day: %v, Turn: %v, cid: %v", v.Verse, v.TimeZone, v.Day, v.TurnInDay, v.IpfsCid)

	leagueProcessor := NewLeagueProcessor(contracts, useractionsPublisher)
	if err := leagueProcessor.Process(tx, v); err != nil {
		return err
	}
	u, err := universe.NewFromStorage(tx, int(v.TimeZone))
	if err != nil {
		return err
	}
	universeHash, err := u.Hash()
	if err != nil {
		return err
	}

	verse := storage.Verse{}
	verse.VerseNumber = v.Verse.Int64()
	verse.Root = hex.EncodeToString(universeHash[:])
	if err := verse.Insert(tx); err != nil {
		return err
	}

	if staker != nil {
		if err := staker.Play(*contracts, v.Verse.Int64(), universeHash); err != nil {
			return err
		}
	}
	return nil
}
