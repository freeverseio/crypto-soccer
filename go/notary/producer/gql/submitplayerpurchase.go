package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) SubmitPlayerPurchase(args struct {
	Input input.SubmitPlayerPurchaseInput
}) (graphql.ID, error) {
	log.Debugf("GeneratePlayerIDs %v", args)

	return graphql.ID(""), errors.New("not implemented")
}
