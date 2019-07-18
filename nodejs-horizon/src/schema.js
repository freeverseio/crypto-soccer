const { gql } = require('apollo-server');

const typeDefs = gql`
  type Query {
    ping: Boolean!
  }

  type Mutation {
    createTeam(name: String!, owner: String!): Boolean,
    createLeague(
      initBlock: Int!
      step: Int!
      teamIds: [ID!]!
      tactics: [[Int!]!]
    ): Boolean
  }`;

module.exports = typeDefs;