package relay

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	//"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/relay/contracts/updates"
	"github.com/freeverseio/crypto-soccer/relay/storage"
)

type Processor struct {
	client        *ethclient.Client
	privateKey    *ecdsa.PrivateKey
	publicAddress common.Address
	db            *storage.Storage
	updates       *updates.Updates
	count         int64
}

// *****************************************************************************
// public
// *****************************************************************************

func NewProcessor(
	client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	db *storage.Storage,
	updates *updates.Updates,
) (*Processor, error) {

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error obtaining publicKey")
	}

	publicAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &Processor{
		client,
		privateKey,
		publicAddress,
		db,
		updates,
		0,
	}, nil
}

func (p *Processor) Process() error {
	return p.computeActionsRoot()
}

// *****************************************************************************
// private
// *****************************************************************************
func (p *Processor) computeActionsRoot() error {
	nonce, err := p.client.PendingNonceAt(context.Background(), p.publicAddress)
	if err != nil {
		return err
	}

	gasPrice, err := p.client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	auth := bind.NewKeyedTransactor(p.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice

	session := updates.UpdatesSession{
		p.updates,
		bind.CallOpts{},
		*auth,
	}

	//currentVerse, err := updates.currentVerse(nil)

	tactic := &storage.Tactic{}
	tactic.TeamID = big.NewInt(p.count)
	p.count += 1
	root, err := tactic.Hash()
	if err != nil {
		return err
	}
	log.Infof("[relay] submitActionsRoot %v", root)
	_, err = session.SubmitActionsRoot(root)
	return err
}
