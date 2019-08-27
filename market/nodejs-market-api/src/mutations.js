const { makeExtendSchemaPlugin, gql } = require("graphile-utils");

const MyPlugin = makeExtendSchemaPlugin(build => {
  // Get any helpers we need from `build`
  const { pgSql: sql, inflection } = build;

  return {
    typeDefs: gql`
      input PlayerSaleOrderInput {
        playerId: String!
      }

      extend type Mutation {
        createPlayerSaleOrder(input: PlayerSaleOrderInput!): Boolean
      }
    `,
    resolvers: {
      Mutation: {
        createPlayerSaleOrder: (_, { input }) =>  {
          return true;
        }
      }
    }
  };
});

module.exports = MyPlugin;