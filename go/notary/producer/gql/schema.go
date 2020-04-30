//go:generate go-bindata -ignore=\.go -pkg=schema -o=bindata.go ./...
package gql

const Schema = ` 
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

	input GetWorldPlayersInput {
		signature: String!
		teamId: ID!
	}

	input SubmitPlayerPurchaseInput {
		signature: String!
		purchaseId: ID!
		playerId: ID!
		teamId: ID!
	}

	type Query {
		getWorldPlayers(input: GetWorldPlayersInput!): [ID!]! 
	}

	type Mutation {
        createAuction(input: CreateAuctionInput!): ID!
        cancelAuction(input: CancelAuctionInput!): ID!
		createBid(input: CreateBidInput!): ID!
		submitPlayerPurchase(input: SubmitPlayerPurchaseInput!): ID!
	}
`
