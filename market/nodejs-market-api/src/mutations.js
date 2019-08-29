const { makeExtendSchemaPlugin, gql } = require("graphile-utils");

const MyPlugin = makeExtendSchemaPlugin(build => {
  // Get any helpers we need from `build`
  const { pgSql: sql, inflection } = build;

  return {
    typeDefs: gql`
       extend type Mutation {
        createPlayerSaleOrder(input: PlayerSaleOrderInput!): ID
        deletePlayerSaleOrder(playerId: ID!): ID
        createPlayerBuyOrder(input: PlayerBuyOrderInput!): ID
        deletePlayerBuyOrder(playerId: ID!): ID
      }
    `,
    resolvers: {
      Mutation: {
        createPlayerSaleOrder: async (_, { input }, context) =>  {
          const { playerId } = input;
          const query = sql.query`INSERT INTO playerSaleOrders (playerId) VALUES (${sql.value(playerId)})`;
          const {text, values} = sql.compile(query);
          await context.pgClient.query(text, values);
          return playerId;
        },
        deletePlayerSaleOrder: async (_, {playerId}, context) => {
          const query = sql.query`DELETE FROM playerSaleOrders WHERE playerId=${sql.value(playerId)}`;
          const {text, values} = sql.compile(query);
          await context.pgClient.query(text, values);
          return playerId;
        },
        createPlayerBuyOrder: async (_, {input}, context) => {
          const { playerId } = input;
          const query = sql.query`INSERT INTO playerBuyOrders (playerId) VALUES (${sql.value(playerId)})`;
          const {text, values} = sql.compile(query);
          await context.pgClient.query(text, values);
          return playerId;
        },
        deletePlayerBuyOrder: async (_, {playerId}, context) => {
          const query = sql.query`DELETE FROM playerBuyOrders WHERE playerId=${sql.value(playerId)}`;
          const {text, values} = sql.compile(query);
          await context.pgClient.query(text, values);
          return playerId;
        },
      }
    }
  };
});

module.exports = MyPlugin;