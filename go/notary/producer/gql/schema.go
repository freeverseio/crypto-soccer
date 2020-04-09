//go:generate go-bindata -ignore=\.go -pkg=schema -o=bindata.go ./...
package gql

const Schema = ` 
	type Query {
		ping: Boolean!,
	}

	input AuctionInput {
  		playerId: String!
  		currencyId: Int!
  		price: Int!
  		rnd: Int!
  		validUntil: String!
  		signature: String!
	}

	input BidInput {
		auction: ID!
  		extraPrice: Int!
  		rnd: Int!
  		teamId: String!
  		signature: String!
	}

	scalar UUID 

	type Mutation {
        createAuction(input: AuctionInput!): UUID!
        // deleteAuction(uuid: UUID!): Boolean!
        // createBid(input: BidInput!): Boolean!
	}
`
