package gql

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/names"
)

type Resolver struct {
	ch                chan interface{}
	contracts         contracts.Contracts
	namesdb           *names.Generator
	googleCredentials []byte
	db                *sql.DB
}

func NewResolver(
	ch chan interface{},
	contracts contracts.Contracts,
	namesdb *names.Generator,
	googleCredentials []byte,
	db *sql.DB,
) *Resolver {
	resolver := Resolver{}
	resolver.ch = ch
	resolver.contracts = contracts
	resolver.namesdb = namesdb
	resolver.googleCredentials = googleCredentials
	resolver.db = db
	return &resolver
}
