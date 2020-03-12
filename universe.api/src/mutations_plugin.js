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
        ): Boolean,
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