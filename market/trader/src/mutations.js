const { makeExtendSchemaPlugin, gql } = require("graphile-utils");

const MyPlugin = makeExtendSchemaPlugin(build => {
  // Get any helpers we need from `build`
  const { pgSql: sql, inflection } = build;

  return {
    typeDefs: gql`
       extend type Mutation {
        createPlayerSellOrder(input: PlayerSellOrderInput!): BigFloat
        deletePlayerSellOrder(playerId: BigFloat!): BigFloat
        createPlayerBuyOrder(input: PlayerBuyOrderInput!): BigFloat
        deletePlayerBuyOrder(playerId: BigFloat!): BigFloat
      }
    `,
    resolvers: {
      Mutation: {
        createPlayerSellOrder: async (_, { input }, context) =>  {
          const { playerid, currencyid, price, rnd, validuntil, typeoftx, signature } = input;
          const query = sql.query`INSERT INTO player_sell_orders (playerId, currencyId, price, rnd, validUntil, typeOfTx, signature) VALUES (
            ${sql.value(playerid)}, 
            ${sql.value(currencyid)}, 
            ${sql.value(price)},
            ${sql.value(rnd)},
            ${sql.value(validuntil)},
            ${sql.value(typeoftx)},
            ${sql.value(signature)}
            )`;
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
          const { playerid, teamid, signature } = input;
          const query = sql.query`INSERT INTO player_buy_orders (playerId, teamId, signature) VALUES (
            ${sql.value(playerid)}, 
            ${sql.value(teamid)}, 
            ${sql.value(signature)}
            )`;
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