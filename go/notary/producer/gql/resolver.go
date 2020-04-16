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
	resolver := Resolver{}
	resolver.ch = ch
	resolver.contracts = contracts
	resolver.db = db
	return &resolver
}

func (b *Resolver) Ping() bool {
	return true
}
