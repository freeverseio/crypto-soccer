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

	input GeneratePlayerIdsInput {
		signature: String!
		seed: Int!
	}

	input SubmitPlayerPurchaseInput {
		signature: String!
		purchaseId: ID!
		playerId: ID!
		teamId: ID!
	}

	type Mutation {
        createAuction(input: CreateAuctionInput!): ID!
        cancelAuction(input: CancelAuctionInput!): ID!
		createBid(input: CreateBidInput!): ID!
		generatePlayerIds(input: GeneratePlayerIdsInput!): [ID!]! 
		submitPlayerPurchase(input: SubmitPlayerPurchaseInput!): ID!
	}
`
