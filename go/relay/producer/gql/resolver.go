package gql

import "github.com/freeverseio/crypto-soccer/go/contracts"

type Resolver struct {
	ch        chan interface{}
	contracts contracts.Contracts
}

func NewResolver(
	ch chan interface{},
	contracts contracts.Contracts,
) *Resolver {
	return &Resolver{
		ch:        ch,
		contracts: contracts,
	}
}

func (b *Resolver) Ping() bool {
	return true
}
