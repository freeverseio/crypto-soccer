//go:generate go-bindata -ignore=\.go -pkg=schema -o=bindata.go ./...
package gql

const Schema = ` 
	type Query {
		ping: Boolean!,
	}

	input CreateAuctionInput {
  		signature: String!
  		playerId: String!
  		currencyId: Int!
  		price: Int!
  		validUntil: String!
  		rnd: Int!
	}

	input CancelAuctionInput {
  		signature: String!
		id: ID!
	}

	input CreateBidInput {
  		signature: String!
		auction: ID!
  		extraPrice: Int!
  		rnd: Int!
  		teamId: String!
	}

	type Mutation {
        createAuction(input: CreateAuctionInput!): ID!
        cancelAuction(input: CancelAuctionInput!): ID!
        createBid(input: CreateBidInput!): ID!
	}
`
