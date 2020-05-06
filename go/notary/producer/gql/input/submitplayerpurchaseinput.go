package input

import "github.com/graph-gophers/graphql-go"

const GooglePackage := "package"
const GoogleProductID := "productId"

type SubmitPlayerPurchaseInput struct {
	Signature  string
	PurchaseId graphql.ID
	PlayerId   graphql.ID
	TeamId     graphql.ID
}
