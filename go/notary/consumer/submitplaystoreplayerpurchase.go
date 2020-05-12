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

	if err := AcknowledgeProduct(
		googleCredentials,
		string(in.PackageName),
		string(in.ProductId),
		in.PurchaseToken,
		"hello world!",
	); err != nil {
		return fmt.Errorf("CRITIC: order with purchaseToken %v with player %v: %v", in.PurchaseToken, playerId, err.Error())
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

	return nil
}

func AcknowledgeProduct(
	credentials []byte,
	packageName string,
	productID string,
	token string,
	payload string,
) error {
	client, err := playstore.New(credentials)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return client.AcknowledgeProduct(
		ctx,
		packageName,
		productID,
		token,
		payload,
	)
}
