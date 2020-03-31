//go:generate go-bindata -ignore=\.go -pkg=schema -o=bindata.go ./...
package gql

const Schema = ` 
	type Query {
		ping: Boolean!,
	}

	type Mutation {
        transferFirstBotToAddr(
          	timezone: Int!,
          	countryIdxInTimezone: ID!,
          	address: String!
		): Boolean!,
	}
`
