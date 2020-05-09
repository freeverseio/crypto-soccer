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
)

func SubmitPlayStorePlayerPurchase(
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	googleCredentials []byte,
	in input.SubmitPlayStorePlayerPurchaseInput,
) error {
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
	receipt, err := helper.WaitReceipt(contracts.Client, tx, 60)
	if err != nil {
		return err
	}
	if receipt.Status == 0 {
		return err
	}

	ctx := context.Background()
	err = client.AcknowledgeProduct(
		ctx,
		string(in.PackageName),
		string(in.ProductId),
		in.PurchaseToken,
		receipt.TxHash.String(),
	)
	if err != nil {
		return fmt.Errorf("CRITIC: order with purchaseToken %v with player %v: %v", in.PurchaseToken, playerId, err.Error())
	}
	return nil
}
