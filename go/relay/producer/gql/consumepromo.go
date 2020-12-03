package gql

import "github.com/graph-gophers/graphql-go"

type ConsumePromoInput struct {
	Signature string
	PlayerId  string
	TeamId    string
}

func (b *Resolver) ConsumePromo(args struct{ Input ConsumePromoInput }) (graphql.ID, error) {
	return graphql.ID(""), nil
}
