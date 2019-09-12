package processor

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/market/notary/contracts/assets"
	"github.com/freeverseio/crypto-soccer/market/notary/storage"
	log "github.com/sirupsen/logrus"
)

type Processor struct {
	db        *storage.Storage
	client    *ethclient.Client
	assets    *assets.Assets
	freeverse *ecdsa.PrivateKey
}

func NewProcessor(db *storage.Storage, ethereumClient *ethclient.Client, assetsContract *assets.Assets, freeverse *ecdsa.PrivateKey) (*Processor, error) {
	return &Processor{db, ethereumClient, assetsContract, freeverse}, nil
}

func RSV(signature string) (r [32]byte, s [32]byte, v uint8, err error) {
	signature = signature[2:] // remove 0x
	vect, err := hex.DecodeString(signature[0:64])
	if err != nil {
		return r, s, v, err
	}
	copy(r[:], vect)
	vect, err = hex.DecodeString(signature[64:128])
	if err != nil {
		return r, s, v, err
	}
	copy(s[:], vect)
	vect, err = hex.DecodeString(signature[128:130])
	v = vect[0]
	return r, s, v, err
}

func (b *Processor) HashPrivateMsg(currencyId uint8, price *big.Int, rnd *big.Int) ([32]byte, error) {
	privateHash, err := b.assets.HashPrivateMsg(
		&bind.CallOpts{},
		currencyId,
		price,
		rnd,
	)
	return privateHash, err
}

func (b *Processor) HashBuyerMessage(hashPrivateMessage [32]byte, validUntil *big.Int, playerId *big.Int, typeOfTx uint8) ([32]byte, error) {
	hash, err := b.assets.BuildPutForSaleTxMsg(
		&bind.CallOpts{},
		hashPrivateMessage,
		validUntil,
		playerId,
		typeOfTx,
	)
	if err != nil {
		return hash, err
	}
	hash, err = b.assets.Prefixed(&bind.CallOpts{}, hash)
	return hash, err
}

func (b *Processor) Process() error {
	log.Info("Processing")

	orders, err := b.db.GetOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		log.Infof("[broker] player %v -> team %v", order.SellOrder.PlayerId, order.BuyOrder.TeamId)

		privHash, err := b.HashPrivateMsg(
			order.SellOrder.CurrencyId,
			order.SellOrder.Price,
			order.SellOrder.Rnd,
		)
		if err != nil {
			log.Error(err)
		}
		rSeller, sSeller, vSeller, err := RSV(order.SellOrder.Signature)
		rBuyer, sBuyer, vBuyer, err := RSV(order.BuyOrder.Signature)
		var sigs [6][32]byte
		var vs [2]uint8
		sigs[0], err = b.HashBuyerMessage(privHash, order.SellOrder.ValidUntil, order.SellOrder.PlayerId, order.SellOrder.TypeOfTx)
		sigs[1] = rSeller
		sigs[2] = sSeller
		vs[0] = vSeller
		sigs[4] = rBuyer
		sigs[5] = sBuyer
		vs[1] = vBuyer
		_, err = b.assets.FreezePlayer(
			bind.NewKeyedTransactor(b.freeverse),
			privHash,
			order.SellOrder.ValidUntil,
			order.SellOrder.PlayerId,
			order.SellOrder.TypeOfTx,
			order.BuyOrder.TeamId,
			sigs,
			vs,
		)
		if err != nil {
			log.Error(err)
		}
		_, err = b.assets.CompleteFreeze(
			bind.NewKeyedTransactor(b.freeverse),
			order.SellOrder.PlayerId,
		)
		if err != nil {
			log.Error(err)
		}
		// _, err = b.assets.TransferPlayer(bind.NewKeyedTransactor(b.freeverse), playerId, teamId)
		// if err != nil {
		// 	log.Error(err)
		// }
		err = b.db.DeleteOrder(order.SellOrder.PlayerId)
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}
