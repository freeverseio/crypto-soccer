package consumer

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/awa/go-iap/playstore"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/androidpublisher/v3"
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
			contracts,
			pvc,
			googleCredentials,
			iapTestOn,
			order,
		); err != nil {
			return err
		}
		if err := order.UpdateState(tx); err != nil {
			return err
		}
	}

	return nil
}

func processPlaystoreOrder(
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	googleCredentials []byte,
	iapTestOn bool,
	order *storage.PlaystoreOrder,
) error {
	playerId, _ := new(big.Int).SetString(order.PlayerId, 10)
	if playerId == nil {
		setState(order, storage.PlaystoreOrderFailed, "invalid player")
		return nil
	}
	teamId, _ := new(big.Int).SetString(order.TeamId, 10)
	if teamId == nil {
		setState(order, storage.PlaystoreOrderFailed, "invalid team")
		return nil
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
		setState(order, storage.PlaystoreOrderFailed, err.Error())
		return nil
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
			setState(order, storage.PlaystoreOrderFailed, err.Error())
			return nil
		}
		if _, err = helper.WaitReceipt(contracts.Client, tx, 60); err != nil {
			setState(order, storage.PlaystoreOrderFailed, err.Error())
			return nil
		}
		log.Infof("[consumer|iap] orderId %v playerId %v assigned to teamId %v", purchase.OrderId, playerId, teamId)
	}

	if !isTestPurchase(purchase) {
		payload := fmt.Sprintf("playerId: %v", order.PlayerId)
		if err := client.AcknowledgeProduct(
			ctx,
			order.PackageName,
			order.ProductId,
			order.PurchaseToken,
			payload,
		); err != nil {
			setState(order, storage.PlaystoreOrderFailed, err.Error())
			return err
		}
	}

	setState(order, storage.PlaystoreOrderComplete, "")
	return nil
}

func isTestPurchase(purchase *androidpublisher.ProductPurchase) bool {
	if purchase.PurchaseType != nil {
		if *purchase.PurchaseType == 0 { // Test
			return true
		}
	}
	return false
}

func setState(order *storage.PlaystoreOrder, state storage.PlaystoreOrderState, extra string) {
	if state == storage.PlaystoreOrderFailed {
		log.Warnf("order %v in state %v with %v", order.OrderId, state, extra)
	}
	order.State = state
	order.StateExtra = extra
}
