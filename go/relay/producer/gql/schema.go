//go:generate go-bindata -ignore=\.go -pkg=schema -o=bindata.go ./...
package gql

const Schema = ` 
	input SetTeamNameInput {
		signature: String!
		teamId: ID!
		name: String!
	}

	type Mutation {
        transferFirstBotToAddr(
          	timezone: Int!,
          	countryIdxInTimezone: ID!,
          	address: String!
		): Boolean!
		setTeamName (input: SetTeamNameInput!): ID!
	}

	type Query {
		ping: Boolean!
	}
`
