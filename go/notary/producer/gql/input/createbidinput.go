package input

import "github.com/graph-gophers/graphql-go"

type CreateBidInput struct {
	Signature  string
	Auction    graphql.ID
	ExtraPrice int
	Rnd        int
	TeamId     string
}
