package consumer

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/marketpay"
	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/notary/producer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	ch                chan interface{}
	contracts         contracts.Contracts
	pvc               *ecdsa.PrivateKey
	market            marketpay.MarketPayService
	googleCredentials []byte
	namesdb           *names.Generator
	iapTestOn         bool
	service           storage.StorageService
}

func New(
	ch chan interface{},
	market marketpay.MarketPayService,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	googleCredentials []byte,
	namesdb *names.Generator,
	iapTestOn bool,
	service storage.StorageService,
) (*Consumer, error) {
	consumer := Consumer{}
	consumer.ch = ch
	consumer.contracts = contracts
	consumer.pvc = pvc
	consumer.market = market
	consumer.googleCredentials = googleCredentials
	consumer.namesdb = namesdb
	consumer.iapTestOn = iapTestOn
	consumer.service = service
	return &consumer, nil
}

func (b *Consumer) Consume(event interface{}) error {
	tx, err := b.service.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	switch in := event.(type) {
	case input.CreateAuctionInput:
		log.Debug("Received CreateAuctionInput")
		return CreateAuction(tx, in)
	case input.CancelAuctionInput:
		return CancelAuction(tx, in)
	case input.CreateBidInput:
		log.Debug("Received CreateBidInput")
		return CreateBid(tx, in)
	case producer.ProcessEvent:
		log.Info("[consumer] process auctions")
		return ProcessAuctions(
			tx,
			b.market,
			b.contracts,
			b.pvc,
		)
	case producer.PlaystoreOrderEvent:
		log.Info("[consumer] process playstore events")
		return ProcessPlaystoreOrders(
			tx,
			b.contracts,
			b.pvc,
			b.googleCredentials,
			b.namesdb,
			b.iapTestOn,
		)
	case input.SubmitPlayStorePlayerPurchaseInput:
		log.Debug("Received SubmitPlayStorePlayerPurchaseInput")
		return SubmitPlayStorePlayerPurchase(tx, in)
	case input.DismissPlayerInput:
		log.Debug("Received DismissPlayerInput")
		return DismissPlayer(b.contracts, b.pvc, in)
	case input.CompletePlayerTransitInput:
		log.Debug("Received CompletePlayerTransit")
		return CompletePlayerTransit(b.contracts, b.pvc, in)
	case input.CreateOfferInput:
		log.Debug("Received CreateOfferInput")
		return CreateOffer(tx, in, b.contracts)
	case input.AcceptOfferInput:
		log.Debug("Received CreateAuctionInput")
		return AcceptOffer(tx, in)
	case input.CancelOfferInput:
		log.Debug("Received CancelOfferInput")
		return CancelOffer(tx, in)
	case producer.ProcessOfferEvent:
		log.Info("[consumer] process offer to expire")
		return ProcessOffers(tx)
	default:
		return fmt.Errorf("unknown event: %+v", event)
	}
}

func (b *Consumer) Start() {
	for {
		event := <-b.ch
		if err := b.Consume(event); err != nil {
			log.Error(err)
		}
	}
}
