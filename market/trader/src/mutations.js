const { makeExtendSchemaPlugin, gql } = require("graphile-utils");

const MyPlugin = makeExtendSchemaPlugin(build => {
  // Get any helpers we need from `build`
  const { pgSql: sql, inflection } = build;

  return {
    typeDefs: gql`
       extend type Mutation {
        createAuction(input: AuctionInput!): Boolean
      }
    `,
    resolvers: {
      Mutation: {
        createAuction: async (_, { input }, context) =>  {
          const { uuid, playerId, currencyId, price, rnd, validUntil, signature } = input;
          const query = sql.query`INSERT INTO auctions (uuid, player_id, currency_id, price, rnd, valid_until, signature) VALUES (
            ${sql.value(uuid)},
            ${sql.value(playerId)}, 
            ${sql.value(currencyId)}, 
            ${sql.value(price)},
            ${sql.value(rnd)},
            ${sql.value(validUntil)},
            ${sql.value(signature)}
            )`;
          const {text, values} = sql.compile(query);
          await context.pgClient.query(text, values);
          return true;
        },
        // deletePlayerSellOrder: async (_, {playerId}, context) => {
        //   const query = sql.query`DELETE FROM player_sell_orders WHERE playerId=${sql.value(playerId)}`;
        //   const {text, values} = sql.compile(query);
        //   await context.pgClient.query(text, values);
        //   return playerId;
        // },
        // createPlayerBuyOrder: async (_, {input}, context) => {
        //   const { playerid, teamid, signature } = input;
        //   const query = sql.query`INSERT INTO player_buy_orders (playerId, teamId, signature) VALUES (
        //     ${sql.value(playerid)}, 
        //     ${sql.value(teamid)}, 
        //     ${sql.value(signature)}
        //     )`;
        //   const {text, values} = sql.compile(query);
        //   await context.pgClient.query(text, values);
        //   return playerid;
        // },
        // deletePlayerBuyOrder: async (_, {playerId}, context) => {
        //   const query = sql.query`DELETE FROM player_buy_orders WHERE playerId=${sql.value(playerId)}`;
        //   const {text, values} = sql.compile(query);
        //   await context.pgClient.query(text, values);
        //   return playerId;
        // },
      }
    }
  };
});

module.exports = MyPlugin;