//go:generate go-bindata -ignore=\.go -pkg=schema -o=bindata.go ./...
package gql

const Schema = ` 
	input ConsumePromoInput {
		signature: String!
		playerId: ID!
		teamId: ID!
	}
	
	type Mutation {
        transferFirstBotToAddr(
          	timezone: Int!,
          	countryIdxInTimezone: ID!,
          	address: String!
		): Boolean!
		consumePromo(input: ConsumePromoInput!): ID!
	}

	type Query {
		ping: Boolean!
	}
`
