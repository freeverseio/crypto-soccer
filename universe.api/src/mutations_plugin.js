const { makeExtendSchemaPlugin, gql } = require("graphile-utils");
const Resolvers = require("./resolvers");

const MyPlugin = (assets, from) => {
  return makeExtendSchemaPlugin(build => {
    // Get any helpers we need from `build`
    const { pgSql: sql, inflection } = build;

    return {
      typeDefs: gql`
      extend type Mutation {
        transferFirstBotToAddr(
          timezone: Int,
          countryIdxInTimezone: ID!,
          address: String!
        ): Boolean
        setTactic(
          teamId: String!
          tacticId: Int!
          shirt0: Int!
          shirt1: Int!
          shirt2: Int!
          shirt3: Int!
          shirt4: Int!
          shirt5: Int!
          shirt6: Int!
          shirt7: Int!
          shirt8: Int!
          shirt9: Int!
          shirt10: Int!
          shirt11: Int!
          shirt12: Int!
          shirt13: Int!
          extraAttack1: Boolean!
          extraAttack2: Boolean!
          extraAttack3: Boolean!
          extraAttack4: Boolean!
          extraAttack5: Boolean!
          extraAttack6: Boolean!
          extraAttack7: Boolean!
          extraAttack8: Boolean!
          extraAttack9: Boolean!
          extraAttack10: Boolean!
        ): Boolean
        setTraining(
          teamId: String!
          specialPlayerShirt: Int!
          goalkeepersDefence: Int!
          goalkeepersSpeed: Int!
          goalkeepersPass: Int!
          goalkeepersShoot: Int!
          goalkeepersEndurance: Int!
          defendersDefence: Int!
          defendersSpeed: Int!
          defendersPass: Int!
          defendersShoot: Int!
          defendersEndurance: Int!
          midfieldersDefence: Int!
          midfieldersSpeed: Int!
          midfieldersPass: Int!
          midfieldersShoot: Int!
          midfieldersEndurance: Int!
          attackersDefence: Int!
          attackersSpeed: Int!
          attackersPass: Int!
          attackersShoot: Int!
          attackersEndurance: Int!
          specialPlayerDefence: Int!
          specialPlayerSpeed: Int!
          specialPlayerPass: Int!
          specialPlayerShoot: Int!
          specialPlayerEndurance: Int!
        ): Boolean
        createSpecialPlayer(
          playerId: String!,
          name: String!,
          defence: Int!,
          speed: Int!,
          pass: Int!,
          shoot: Int!,
          endurance: Int!,
          preferredPosition: String!,
          potential: Int!,
          dayOfBirth: Int!
        ): Boolean,
        deleteSpecialPlayer(
          playerId: String!
        ): Boolean
      }`,
      resolvers: Resolvers(sql, assets, from),
    }
  });
};

module.exports = MyPlugin;