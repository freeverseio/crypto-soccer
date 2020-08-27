package consumer

import (
	"database/sql"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/complementarydata"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/relay/producer"
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/useractions"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	ch                        chan interface{}
	client                    *ethclient.Client
	auth                      *bind.TransactOpts
	contracts                 contracts.Contracts
	useractionsPublishService useractions.UserActionsPublishService
	db                        *sql.DB
	complementaryData         complementarydata.ComplementaryData
}

func NewConsumer(
	ch chan interface{},
	client *ethclient.Client,
	auth *bind.TransactOpts,
	contracts contracts.Contracts,
	useractionsPublishService useractions.UserActionsPublishService,
	db *sql.DB,
) *Consumer {
	return &Consumer{
		ch,
		client,
		auth,
		contracts,
		useractionsPublishService,
		db,
		complementarydata.ComplementaryData{},
	}
}

func (b *Consumer) Start() {
	firstBotTransfer := NewFirstBotTransfer(b.client, b.auth, b.contracts)
	actionsSubmitter := NewActionsSubmitter(b.client, b.auth, b.contracts, b.useractionsPublishService)
	for {
		event := <-b.ch
		switch ev := event.(type) {
		case gql.TransferFirstBotToAddrInput:
			log.Infof("[relay|consumer] Trasfer First Bot TX: %v Country: %v to %v", ev.Timezone, ev.CountryIdxInTimezone, ev.Address)
			if err := firstBotTransfer.Process(ev); err != nil {
				log.Error(err)
			}
		case producer.SubmitActionsEvent:
			log.Debug("[relay|consumer] Relay sumbit action event")
			tx, err := b.db.Begin()
			if err != nil {
				log.Error(err)
				break
			}
			if err = actionsSubmitter.Process(tx); err != nil {
				tx.Rollback()
				log.Error(err)
				break
			}
			if err = tx.Commit(); err != nil {
				log.Error(err)
			}
		default:
			log.Errorf("unknown event: %v", event)
		}
	}
}
