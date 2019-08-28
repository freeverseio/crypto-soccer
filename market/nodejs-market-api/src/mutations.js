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
        createPlayerSaleOrder(input: PlayerSaleOrderInput!): Boolean!
      }
    `,
    resolvers: {
      Mutation: {
        createPlayerSaleOrder: async (_, { input }, context) =>  {
          const { playerId } = input;
          const query = sql.query`INSERT INTO playerSaleOrders (playerId) VALUES (${sql.value(playerId)})`;
          const {text, values} = sql.compile(query);
          await context.pgClient.query(text, values);
          return true;
        }
      }
    }
  };
});

module.exports = MyPlugin;