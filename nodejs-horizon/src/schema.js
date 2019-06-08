const { gql } = require('apollo-server');

const typeDefs = gql`
  type Query {
    countTeams: String!
    getTeam(id: ID!): Team
    allTeams: [Team]!
    getPlayer(id: ID!): Player
    countLeagues: String!
  }

  type Mutation {
    createTeam(name: String!, owner: String!): Team,
    createLeague(
      initBlock: Int!
      step: Int!
      teamIds: [ID!]!
      tactics: [[Int!]!]
    ): League
  }

  type Subscription {
    teamCreated: ID!
  }

  type Team {
    id: ID!
    name: String!
    players: [Player!]
  }

  type League {
    id: ID!
    initBlock: Int!
    step: Int!
    nTeams: Int!
  }
  
  type Player {
    id: ID!
    name: String!
    defence: Int!
    speed: Int!
    pass: Int!
    shoot: Int!
    endurance: Int!
    team: Team!
  }`;

module.exports = typeDefs;