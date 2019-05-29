const { gql } = require('apollo-server');

const typeDefs = gql`
  type Query {
    settings: Settings!
    countTeams: String!
    teamById(id: ID!): Team
    allTeams: [Team]
  }

  type Mutation {
    createTeam(name: String!, owner: String!): String
  }

  type Subscription {
    teamCreated: ID!
  }

  type Settings {
    network_id: String
    assetsContractAddress: String
  }

  type Team {
    id: ID!
    name: String!
    playerIds: [ID!]
  }`;

module.exports = typeDefs;