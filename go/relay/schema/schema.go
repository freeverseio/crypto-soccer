//go:generate go-bindata -ignore=\.go -pkg=schema -o=bindata.go ./...
package schema

const Schema = ` 
	type Mutation {
        transferFirstBotToAddr(
          	timezone: Int,
          	countryIdxInTimezone: ID!,
          	address: String!
        ): Boolean,
`
