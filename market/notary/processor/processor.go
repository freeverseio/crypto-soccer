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

func (b *Processor) HashSellMessage(currencyId uint8, price *big.Int, rnd *big.Int, validUntil *big.Int, playerId *big.Int, typeOfTx uint8) ([32]byte, error) {
	var hash [32]byte
	hashPrivateMessage, err := b.assets.HashPrivateMsg(
		&bind.CallOpts{},
		currencyId,
		price,
		rnd,
	)
	if err != nil {
		return hash, err
	}
	hash, err = b.assets.BuildPutForSaleTxMsg(
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

func (b *Processor) HashBuyMessage(currencyId uint8, price *big.Int, rnd *big.Int, validUntil *big.Int, playerId *big.Int, typeOfTx uint8, teamId *big.Int) ([32]byte, error) {
	var hash [32]byte
	hashPrivateMessage, err := b.assets.HashPrivateMsg(
		&bind.CallOpts{},
		currencyId,
		price,
		rnd,
	)
	if err != nil {
		return hash, err
	}
	sellMsgHash, err := b.assets.BuildPutForSaleTxMsg(
		&bind.CallOpts{},
		hashPrivateMessage,
		validUntil,
		playerId,
		typeOfTx,
	)
	if err != nil {
		return hash, err
	}
	prefixedHash, err := b.assets.Prefixed(&bind.CallOpts{}, sellMsgHash)
	if err != nil {
		return hash, err
	}
	hash, err = b.assets.BuildAgreeToBuyTxMsg(
		&bind.CallOpts{},
		prefixedHash,
		teamId,
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

		log.Infof("(1) generate hash private msg")
		privHash, err := b.HashPrivateMsg(
			order.SellOrder.CurrencyId,
			order.SellOrder.Price,
			order.SellOrder.Rnd,
		)
		if err != nil {
			log.Error(err)
			continue
		}

		log.Infof("(2) generate hash sell message")
		var sigs [6][32]byte
		var vs [2]uint8
		sigs[0], err = b.HashSellMessage(order.SellOrder.CurrencyId, order.SellOrder.Price, order.SellOrder.Rnd, order.SellOrder.ValidUntil, order.SellOrder.PlayerId, order.SellOrder.TypeOfTx)
		if err != nil {
			log.Error(err)
			continue
		}
		sigs[1], sigs[2], vs[0], err = RSV(order.SellOrder.Signature)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("(3) generate hash buy message")
		sigs[3], err = b.HashBuyMessage(order.SellOrder.CurrencyId, order.SellOrder.Price, order.SellOrder.Rnd, order.SellOrder.ValidUntil, order.SellOrder.PlayerId, order.SellOrder.TypeOfTx, order.BuyOrder.TeamId)
		if err != nil {
			log.Error(err)
			continue
		}
		sigs[4], sigs[5], vs[1], err = RSV(order.BuyOrder.Signature)
		if err != nil {
			log.Error(err)
			continue
		}

		log.Infof("(4) freeze player")
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
			continue
		}
		log.Infof("(5) complete freeze")
		_, err = b.assets.CompleteFreeze(
			bind.NewKeyedTransactor(b.freeverse),
			order.SellOrder.PlayerId,
		)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("(6) delete order")
		err = b.db.DeleteOrder(order.SellOrder.PlayerId)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	return nil
}
