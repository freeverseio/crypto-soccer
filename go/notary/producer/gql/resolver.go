package gql

import (
	"database/sql"
	"errors"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/names"
	log "github.com/sirupsen/logrus"
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

func (b *Resolver) push(event interface{}) error {
	select {
	case b.ch <- event:
	default:
		log.Warning("channel is full")
		return errors.New("channel is full")
	}
	return nil
}
