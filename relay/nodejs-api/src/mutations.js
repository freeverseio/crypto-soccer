const { makeExtendSchemaPlugin, gql } = require("graphile-utils");

const MyPlugin = makeExtendSchemaPlugin(build => {
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
      }`,
    resolvers: {
      Mutation: {
        transferFirstBotToAddr: async (_, { timezone, countryIdxInTimezone, address }) => {
          const gas = await assets.methods
            .transferFirstBotToAddr(timezone, countryIdxInTimezone, address)
            .estimateGas();
          await assets.methods
            .transferFirstBotToAddr(timezone, countryIdxInTimezone, address)
            .send({ from, gas });
          return true;
        },
      },
    },
  };
});

module.exports = MyPlugin;