const { makeExtendSchemaPlugin, gql } = require("graphile-utils");

const MyPlugin = makeExtendSchemaPlugin(build => {
  // Get any helpers we need from `build`
  const { pgSql: sql, inflection } = build;

  return {
    typeDefs: gql`
       extend type Mutation {
        createAuction(input: AuctionInput!): Boolean
        deleteAuction(uuid: UUID!): Boolean
        createBid(input: BidInput!): Boolean
      }
    `,
    resolvers: {
      Mutation: {
        createAuction: async (_, { input }, context) =>  {
          const { uuid, playerId, currencyId, price, rnd, validUntil, signature } = input;
          const query = sql.query`INSERT INTO auctions (uuid, player_id, currency_id, price, rnd, valid_until, signature, state) VALUES (
            ${sql.value(uuid)},
            ${sql.value(playerId)}, 
            ${sql.value(currencyId)}, 
            ${sql.value(price)},
            ${sql.value(rnd)},
            ${sql.value(validUntil)},
            ${sql.value(signature)},
            ${sql.value('STARTED')}
            )`;
          const {text, values} = sql.compile(query);
          await context.pgClient.query(text, values);
          return true;// TODO return something with sense
        },
        deleteAuction: async (_, {uuid}, context) => {
          const query = sql.query`UPDATE auctions SET state='CANCELLED_BY_SELLER' WHERE uuid=${sql.value(uuid)}`;
          const {text, values} = sql.compile(query);
          await context.pgClient.query(text, values);
          return true; // TODO return something with sense
        },
        createBid: async (_, {input}, context) => {
          const { auction, extraPrice, rnd, teamId, signature } = input;
          const query = sql.query`INSERT INTO bids (auction, extra_price, rnd, team_id, signature, state) VALUES (
            ${sql.value(auction)}, 
            ${sql.value(extraPrice)}, 
            ${sql.value(rnd)},
            ${sql.value(teamId)},
            ${sql.value(signature)},
            ${sql.value('FILED')}
            )`;
          const {text, values} = sql.compile(query);
          await context.pgClient.query(text, values);
          return true;
        },
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