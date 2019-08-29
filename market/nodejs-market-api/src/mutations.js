const { makeExtendSchemaPlugin, gql } = require("graphile-utils");

const MyPlugin = makeExtendSchemaPlugin(build => {
  // Get any helpers we need from `build`
  const { pgSql: sql, inflection } = build;

  return {
    typeDefs: gql`
       extend type Mutation {
        createPlayerSellOrder(input: PlayerSellOrderInput!): ID
        deletePlayerSellOrder(playerId: ID!): ID
        createPlayerBuyOrder(input: PlayerBuyOrderInput!): ID
        deletePlayerBuyOrder(playerId: ID!): ID
      }
    `,
    resolvers: {
      Mutation: {
        createPlayerSellOrder: async (_, { input }, context) =>  {
          const { playerid, price } = input;
          const query = sql.query`INSERT INTO player_sell_orders (playerId, price) VALUES (${sql.value(playerid)}, ${sql.value(price)})`;
          const {text, values} = sql.compile(query);
          await context.pgClient.query(text, values);
          return playerid;
        },
        deletePlayerSellOrder: async (_, {playerId}, context) => {
          const query = sql.query`DELETE FROM player_sell_orders WHERE playerId=${sql.value(playerId)}`;
          const {text, values} = sql.compile(query);
          await context.pgClient.query(text, values);
          return playerId;
        },
        createPlayerBuyOrder: async (_, {input}, context) => {
          const { playerid, price } = input;
          const query = sql.query`INSERT INTO player_buy_orders (playerId, price) VALUES (${sql.value(playerid)}, ${sql.value(price)})`;
          const {text, values} = sql.compile(query);
          await context.pgClient.query(text, values);
          return playerid;
        },
        deletePlayerBuyOrder: async (_, {playerId}, context) => {
          const query = sql.query`DELETE FROM player_buy_orders WHERE playerId=${sql.value(playerId)}`;
          const {text, values} = sql.compile(query);
          await context.pgClient.query(text, values);
          return playerId;
        },
      }
    }
  };
});

module.exports = MyPlugin;