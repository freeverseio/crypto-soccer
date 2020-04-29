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
		auctionId: ID!
	}

	input CreateBidInput {
  		signature: String!
		auctionId: ID!
  		extraPrice: Int!
  		rnd: Int!
  		teamId: String!
	}

	input GeneratePlayerIDsInput {

	}

	input SubmitPlayerPurchaseInput {

	}

	type Mutation {
        createAuction(input: CreateAuctionInput!): ID!
        cancelAuction(input: CancelAuctionInput!): ID!
		createBid(input: CreateBidInput!): ID!
		generatePlayerIDs(input: GeneratePlayerIDsInput!): [ID!]! 
		submitPlayerPurchase(input: SubmitPlayerPurchaseInput!): ID!
	}
`
