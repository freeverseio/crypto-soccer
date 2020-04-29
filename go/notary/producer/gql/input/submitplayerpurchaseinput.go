package input

import "github.com/graph-gophers/graphql-go"

type SubmitPlayerPurchaseInput struct {
	Signature  string
	PurchaseId graphql.ID
	PlayerId   graphql.ID
	TeamId     graphql.ID
}
