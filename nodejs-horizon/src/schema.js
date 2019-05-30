const { gql } = require('apollo-server');

const typeDefs = gql`
  type Query {
    countTeams: String!
    getTeam(id: ID!): Team
    allTeams: [Team]
    getPlayer(id: ID!): Player
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
    players: [Player!]
  }
  
  type Player {
    id: ID!
    name: String!
    defence: Int!
  }`;

module.exports = typeDefs;