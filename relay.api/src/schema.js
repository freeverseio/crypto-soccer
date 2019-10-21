const { gql } = require('apollo-server');

const typeDefs = gql`
  type Query {
    ping: Boolean!
  }

  type Mutation {
    transferFirstBotToAddr(
      timezone: Int,
      countryIdxInTimezone: ID!,
      address: String!
    ): Boolean
  }`;

module.exports = typeDefs;