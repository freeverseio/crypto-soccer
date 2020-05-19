package consumer

import (
	"crypto/ecdsa"
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"google.golang.org/api/androidpublisher/v3"

	log "github.com/sirupsen/logrus"
)

func SubmitPlayStorePlayerPurchase(
	contracts contracts.Contracts,
	tx *sql.Tx,
	pvc *ecdsa.PrivateKey,
	googleCredentials []byte,
	in input.SubmitPlayStorePlayerPurchaseInput,
	iapTestOn bool,
) error {
	log.Debugf("SubmitPlayStorePlayerPurchase %+v", in)

	// playerId, _ := new(big.Int).SetString(string(in.PlayerId), 10)
	// if playerId == nil {
	// 	return fmt.Errorf("invalid playerId %v", in.PlayerId)
	// }
	// teamId, _ := new(big.Int).SetString(string(in.TeamId), 10)
	// if teamId == nil {
	// 	return fmt.Errorf("invalid teamId %v", in.TeamId)
	// }

	data, err := playstore.InappPurchaseDataFromReceipt(in.Receipt)
	if err != nil {
		return err
	}

	order := storage.NewPlaystoreOrder()
	order.OrderId = data.OrderId
	order.PackageName = data.PackageName
	order.ProductId = data.ProductId
	order.PurchaseToken = data.PurchaseToken
	order.PlayerId = string(in.PlayerId)
	order.TeamId = string(in.TeamId)
	order.Signature = in.Signature
	if err := order.Insert(tx); err != nil {
		return err
	}

	// client, err := playstore.New(googleCredentials)
	// if err != nil {
	// 	return err
	// }
	// ctx := context.Background()

	// purchase, err := client.VerifyProduct(
	// 	ctx,
	// 	string(in.PackageName),
	// 	string(in.ProductId),
	// 	in.PurchaseToken,
	// )
	// if err != nil {
	// 	return err
	// }

	// if isTestPurchase(purchase) && !iapTestOn {
	// 	log.Warningf("[consumer|iap] received test orderId %v ... skip creating player", purchase.OrderId)
	// } else {
	// 	auth := bind.NewKeyedTransactor(pvc)
	// 	auth.GasPrice = big.NewInt(1000000000) // in xdai is fixe to 1 GWei
	// 	tx, err := contracts.Market.TransferBuyNowPlayer(
	// 		auth,
	// 		playerId,
	// 		teamId,
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if _, err = helper.WaitReceipt(contracts.Client, tx, 60); err != nil {
	// 		return err
	// 	}
	// 	log.Infof("[consumer|iap] orderId %v playerId %v assigned to teamId %v", purchase.OrderId, playerId, teamId)
	// }

	// payload := fmt.Sprintf("playerId: %v", in.PlayerId)
	// if err := client.AcknowledgeProduct(
	// 	ctx,
	// 	string(in.PackageName),
	// 	string(in.ProductId),
	// 	in.PurchaseToken,
	// 	payload,
	// ); err != nil {
	// 	return err
	// }

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
