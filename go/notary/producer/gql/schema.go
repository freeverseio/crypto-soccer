//go:generate go-bindata -ignore=\.go -pkg=schema -o=bindata.go ./...
package gql

const Schema = ` 
	type Query {
		ping: Boolean!,
	}

	scalar Sign	

	input CreateAuctionInput {
  		signature: Sign!
  		playerId: String!
  		currencyId: Int!
  		price: Int!
  		validUntil: String!
  		rnd: Int!
	}

	input CancelAuctionInput {
		signature: Sign!
	}

	input CreateBidInput {
  		signature: Sign!
		auction: Sign!
  		extraPrice: Int!
  		rnd: Int!
  		teamId: String!
	}

	type Mutation {
        createAuction(input: CreateAuctionInput!): Sign!
        cancelAuction(input: CancelAuctionInput!): Sign!
        createBid(input: CreateBidInput!): Sign!
	}
`
