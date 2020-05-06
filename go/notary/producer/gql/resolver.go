package gql

import (
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/names"
)

type Resolver struct {
	ch        chan interface{}
	contracts contracts.Contracts
	namesdb   *names.Generator
}

func NewResolver(
	ch chan interface{},
	contracts contracts.Contracts,
	namesdb *names.Generator,
) *Resolver {
	resolver := Resolver{}
	resolver.ch = ch
	resolver.contracts = contracts
	resolver.namesdb = namesdb
	return &resolver
}
