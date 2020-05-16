package consumer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/awa/go-iap/playstore"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"google.golang.org/api/androidpublisher/v3"

	log "github.com/sirupsen/logrus"
)

func SubmitPlayStorePlayerPurchase(
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	googleCredentials []byte,
	in input.SubmitPlayStorePlayerPurchaseInput,
	iapTestOn bool,
) error {
	log.Debugf("SubmitPlayStorePlayerPurchase %+v", in)

	playerId, _ := new(big.Int).SetString(string(in.PlayerId), 10)
	if playerId == nil {
		return fmt.Errorf("invalid playerId %v", in.PlayerId)
	}
	teamId, _ := new(big.Int).SetString(string(in.TeamId), 10)
	if teamId == nil {
		return fmt.Errorf("invalid teamId %v", in.TeamId)
	}

	client, err := playstore.New(googleCredentials)
	if err != nil {
		return err
	}
	ctx := context.Background()

	purchase, err := client.VerifyProduct(
		ctx,
		string(in.PackageName),
		string(in.ProductId),
		in.PurchaseToken,
	)
	if err != nil {
		return err
	}

	if isTestPurchase(purchase) && !iapTestOn {
		log.Warningf("[consumer|iap] received test orderId %v ... skip", purchase.OrderId)
		return nil
	} else if isTestPurchase(purchase) && iapTestOn {
		log.Infof("[consumer|iap] test order: skip sending acknowledge to google service")
	} else {
		payload := fmt.Sprintf("playerId: %v", in.PlayerId)
		if err := client.AcknowledgeProduct(
			ctx,
			string(in.PackageName),
			string(in.ProductId),
			in.PurchaseToken,
			payload,
		); err != nil {
			return err
		}
	}

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
