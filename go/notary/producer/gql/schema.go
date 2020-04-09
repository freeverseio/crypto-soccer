//go:generate go-bindata -ignore=\.go -pkg=schema -o=bindata.go ./...
package gql

const Schema = ` 
	type Query {
		ping: Boolean!,
	}

	scalar UUID

	input AuctionInput {
		uuid: UUID!
  		playerId: String!
  		currencyId: Int!
  		price: Int!
  		rnd: Int!
  		validUntil: String!
  		signature: String!
  		stateExtra: String
  		seller: String!
	}

	input BidInput {
		auction: UUID!
  		extraPrice: Int!
  		rnd: Int!
  		teamId: String!
  		signature: String!
  		stateExtra: String
  		paymentId: String
  		paymentDeadline: String
	}

	type Mutation {
        createAuction(input: AuctionInput!): Boolean!
        deleteAuction(uuid: UUID!): Boolean!
        createBid(input: BidInput!): Boolean!
	}
`
