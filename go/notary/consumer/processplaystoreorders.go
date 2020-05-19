package consumer

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"errors"
	"fmt"
	"math/big"

	"github.com/awa/go-iap/playstore"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func ProcessPlaystoreOrders(
	tx *sql.Tx,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	googleCredentials []byte,
	iapTestOn bool,
) error {
	orders, err := storage.PendingPlaystoreOrders(tx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if err := processPlaystoreOrder(
			tx,
			contracts,
			pvc,
			googleCredentials,
			iapTestOn,
			order,
		); err != nil {
			log.Errorf("[consumer|error] %v, playstore order %+v", err.Error(), order)
		}
	}

	return nil
}

func processPlaystoreOrder(
	tx *sql.Tx,
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	googleCredentials []byte,
	iapTestOn bool,
	order storage.PlaystoreOrder,
) error {
	playerId, _ := new(big.Int).SetString(order.PlayerId, 10)
	if playerId == nil {
		return errors.New("invalid playerId")
	}
	teamId, _ := new(big.Int).SetString(order.TeamId, 10)
	if teamId == nil {
		return errors.New("invalid teamId")
	}

	client, err := playstore.New(googleCredentials)
	if err != nil {
		return err
	}
	ctx := context.Background()

	purchase, err := client.VerifyProduct(
		ctx,
		order.PackageName,
		order.ProductId,
		order.PurchaseToken,
	)
	if err != nil {
		return err
	}

	if isTestPurchase(purchase) && !iapTestOn {
		log.Warningf("[consumer|iap] received test orderId %v ... skip creating player", purchase.OrderId)
	} else {
		auth := bind.NewKeyedTransactor(pvc)
		auth.GasPrice = big.NewInt(1000000000) // in xdai is fixe to 1 GWei
		tx, err := contracts.Market.TransferBuyNowPlayer(
			auth,
			playerId,
			teamId,
		)
		if err != nil {
			return err
		}
		if _, err = helper.WaitReceipt(contracts.Client, tx, 60); err != nil {
			return err
		}
		log.Infof("[consumer|iap] orderId %v playerId %v assigned to teamId %v", purchase.OrderId, playerId, teamId)
	}

	payload := fmt.Sprintf("playerId: %v", order.PlayerId)
	if err := client.AcknowledgeProduct(
		ctx,
		order.PackageName,
		order.ProductId,
		order.PurchaseToken,
		payload,
	); err != nil {
		return err
	}
	return nil
}
