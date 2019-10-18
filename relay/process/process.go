package relay

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
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

	rand.Seed(time.Now())

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error obtaining publicKey")
	}

	publicAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &Processor{
		client,
		publicAddress,
		privateKey,
		db,
		updates,
	}, nil
}

func (p *Processor) Process(delta uint64) error {
	return computeActionsRoot()
}

type Action struct {
	value uint64
}

// Hash - computes hash for an Action
func (a *Action) Hash() ([]byte, error) {
	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return computeHash(sha256.New(), data), nil
}

func computeHash(h hash.Hash, data ...[]byte) []byte {
	h.Reset()
	for _, d := range data {
		h.Write(d)
	}
	return h.Sum(nil)
}

// *****************************************************************************
// private
// *****************************************************************************
func (p *Processor) computeActionsRoot() error {
	gasLimit := uint64(21000)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice

	//currentVerse, err := updates.currentVerse(nil)

	root := Action{rand.Uint64()}.Hash()
	log.Infof("[relay] submitActionsRoot %v", root)
	_, err = updates.submitActionsRoot(auth, root)
	return err
}
