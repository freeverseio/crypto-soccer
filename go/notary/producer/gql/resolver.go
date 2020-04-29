package gql

import (
	"github.com/freeverseio/crypto-soccer/go/contracts"
)

type Resolver struct {
	ch        chan interface{}
	contracts contracts.Contracts
}

func NewResolver(
	ch chan interface{},
	contracts contracts.Contracts,
) *Resolver {
	resolver := Resolver{}
	resolver.ch = ch
	resolver.contracts = contracts
	return &resolver
}
