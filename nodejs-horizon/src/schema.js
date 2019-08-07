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
    ): Boolean,
    transferTeam(
      teamId: String!,
      owner: String!
    ): Boolean
  }`;

module.exports = typeDefs;