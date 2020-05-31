package contracts

import (
	"database/sql"
	"errors"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/storage"
)

func NewFromStorage(client *ethclient.Client, tx *sql.Tx) (*Contracts, error) {
	proxyParam, err := storage.ParamByName(tx, ProxyName)
	if err != nil {
		return nil, err
	}
	if proxyParam == nil {
		return nil, errors.New("no proxy address in the storage")
	}

	return NewByProxyAddress(client, proxyParam.Value)
}

func (b Contracts) ToStorage(tx *sql.Tx) error {
	if err := (storage.Param{ProxyName, b.ProxyAddress}).InsertOrUpdate(tx); err != nil {
		return err
	}
	return nil
}
