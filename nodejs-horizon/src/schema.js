const { gql } = require('apollo-server');

const typeDefs = gql`
  type Query {
    ping: Boolean!
  }

  type Mutation {
    createLeague(
      nTeams: Int!
      initBlock: Int!
      step: Int!
    ): Boolean
  }`;

module.exports = typeDefs;