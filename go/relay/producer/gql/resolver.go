package gql

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/contracts"
)

type Resolver struct {
	ch        chan interface{}
	contracts contracts.Contracts
	db        *sql.DB
}

func NewResolver(
	ch chan interface{},
	contracts contracts.Contracts,
	db *sql.DB,
) *Resolver {
	return &Resolver{
		ch:        ch,
		contracts: contracts,
		db:        db,
	}
}

func (b *Resolver) Ping() bool {
	return true
}
