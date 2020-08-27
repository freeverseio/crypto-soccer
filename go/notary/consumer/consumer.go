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
	switch in := event.(type) {
	case input.CreateAuctionInput:
		log.Debug("Received CreateAuctionInput")
		if err := b.service.Begin(); err != nil {
			return err
		}
		if err := CreateAuction(b.service, in); err != nil {
			b.service.Rollback()
			return err
		}
		return b.service.Commit()
	case input.CancelAuctionInput:
		log.Debug("Received CancelAuctionInput")
		if err := b.service.Begin(); err != nil {
			return err
		}
		if err := CancelAuction(b.service, in); err != nil {
			b.service.Rollback()
			return err
		}
		return b.service.Commit()
	case input.CreateBidInput:
		log.Debug("Received CreateBidInput")
		if err := b.service.Begin(); err != nil {
			return err
		}
		if err := CreateBid(b.service, in); err != nil {
			b.service.Rollback()
			return err
		}
		return b.service.Commit()
	case producer.ProcessEvent:
		log.Info("[consumer] process auctions")
		if err := b.service.Begin(); err != nil {
			return err
		}
		if err := ProcessAuctions(
			b.service,
			b.market,
			b.contracts,
			b.pvc,
		); err != nil {
			b.service.Rollback()
			return err
		}
		return b.service.Commit()
	case producer.PlaystoreOrderEvent:
		log.Info("[consumer] process playstore events")
		if err := b.service.Begin(); err != nil {
			return err
		}
		if err := ProcessPlaystoreOrders(
			b.service,
			b.contracts,
			b.pvc,
			b.googleCredentials,
			b.namesdb,
			b.iapTestOn,
		); err != nil {
			b.service.Rollback()
			return err
		}
		return b.service.Commit()
	case input.SubmitPlayStorePlayerPurchaseInput:
		log.Debug("Received SubmitPlayStorePlayerPurchaseInput")
		if err := b.service.Begin(); err != nil {
			return err
		}
		if err := SubmitPlayStorePlayerPurchase(
			b.service,
			in,
		); err != nil {
			b.service.Rollback()
			return err
		}
		return b.service.Commit()
	case input.DismissPlayerInput:
		log.Debug("Received DismissPlayerInput")
		if err := DismissPlayer(
			b.contracts,
			b.pvc,
			in,
		); err != nil {
			return err
		}
	case input.CompletePlayerTransitInput:
		log.Debug("Received CompletePlayerTransit")
		if err := CompletePlayerTransit(b.contracts, b.pvc, in); err != nil {
			return err
		}
	case input.CreateOfferInput:
		log.Debug("Received CreateOfferInput")

		if err := b.service.Begin(); err != nil {
			return err
		}
		if err := CreateOffer(b.service, in, b.contracts); err != nil {
			b.service.Rollback()
			return err
		}
		return b.service.Commit()
	case input.AcceptOfferInput:
		log.Debug("Received CreateAuctionInput")
		if err := b.service.Begin(); err != nil {
			return err
		}
		if err := AcceptOffer(b.service, in); err != nil {
			b.service.Rollback()
			return err
		}
		return b.service.Commit()
	case input.CancelOfferInput:
		log.Debug("Received CancelOfferInput")
		if err := b.service.Begin(); err != nil {
			return err
		}
		if err := CancelOffer(b.service, in); err != nil {
			b.service.Rollback()
			return err
		}
		return b.service.Commit()
	case producer.ProcessOfferEvent:
		log.Info("[consumer] process offer to expire")
		if err := b.service.Begin(); err != nil {
			return err
		}
		if err := ProcessOffers(
			b.service,
		); err != nil {
			b.service.Rollback()
			return err
		}
		return b.service.Commit()
	default:
		return fmt.Errorf("unknown event: %+v", event)
	}
	return nil
}

func (b *Consumer) Start() {
	for {
		event := <-b.ch
		if err := b.Consume(event); err != nil {
			log.Error(err)
		}
	}
}
