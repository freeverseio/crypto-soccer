const { ApolloServer } = require("apollo-server");
const { HttpLink } = require("apollo-link-http");
const { introspectSchema, makeRemoteExecutableSchema, mergeSchemas } = require("graphql-tools");
const fetch = require("node-fetch");
const program = require("commander");
const version = require("../package.json").version;

const createRemoteSchema = async uri => {
  const link = new HttpLink({ uri, fetch });
  const schema = await introspectSchema(link);
  const executableSchema = makeRemoteExecutableSchema({
    schema,
    link
  });
  return executableSchema;
};

program
  .version(version)
  .option("-u, --universeUrl <url>", "graphql universe url", "")
  .option("-m, --marketUrl <url>", "graphql market url", "")
  .option("-r, --relayUrl <url>", "graphql relay url", "")
  .option("-n, --notaryUrl <url>", "graphql notary url", "")
  .option("-c, --enableCors", "enable CORS")
  .parse(process.argv);

const { universeUrl, marketUrl, relayUrl, notaryUrl, enableCors } = program;

console.log("--------------------------------------------------------");
console.log("universeUrl       : ", universeUrl);
console.log("marketUrl         : ", marketUrl);
console.log("relayUrl          : ", relayUrl);
console.log("notaryUrl         : ", notaryUrl);
console.log("--------------------------------------------------------");

const main = async () => {
  const universeRemoteSchema = await createRemoteSchema(universeUrl);
  const marketRemoteSchema = await createRemoteSchema(marketUrl);
  const relayRemoteSchema = (relayUrl !== "") ? await createRemoteSchema(relayUrl) : null;
  const notaryRemoteSchema = (notaryUrl !== "") ? await createRemoteSchema(notaryUrl) : null;

  const linkTypeDefs = `
    extend type Player {
      auctionsByPlayerId(
        first: Int
        last: Int
        offset: Int
        before: Cursor
        after: Cursor
        orderBy: [AuctionsOrderBy!] = [PRIMARY_KEY_ASC]
        condition: AuctionCondition
      ): AuctionsConnection!,
      offersByPlayerId(
        first: Int
        last: Int
        offset: Int
        before: Cursor
        after: Cursor
        orderBy: [OffersOrderBy!] = [PRIMARY_KEY_ASC]
        condition: AuctionCondition
      ): OffersConnection!

    }

    extend type Auction {
      playerByPlayerId: Player
    }

    extend type Bid {
      teamByTeamId: Team
    }

    extend type Offer {
      teamByBuyerTeamId: Team
      playerByPlayerId: Player
    }
  `;

  const resolvers = {
    Player: {
      auctionsByPlayerId: {
        fragment: `... on Player { playerId }`,
        resolve(player, args, context, info) {
          return info.mergeInfo.delegateToSchema({
            schema: marketRemoteSchema,
            operation: 'query',
            fieldName: 'allAuctions',
            args: {
              condition: {
                playerId: player.playerId
              }
            },
            context,
            info,
          })
        }
      },
      offersByPlayerId: {
        fragment: `... on Player { playerId }`,
        resolve(player, args, context, info) {
          return info.mergeInfo.delegateToSchema({
            schema: marketRemoteSchema,
            operation: 'query',
            fieldName: 'allOffers',
            args: {
              condition: {
                playerId: player.playerId
              }
            },
            context,
            info,
          })
        }
      },
    },
    Auction: {
      playerByPlayerId: {
        fragment: `... on Auction { playerId }`,
        resolve(auction, args, context, info) {
          return info.mergeInfo.delegateToSchema({
            schema: universeRemoteSchema,
            operation: 'query',
            fieldName: 'playerByPlayerId',
            args: {
              playerId: auction.playerId,
            },
            context,
            info,
          })
        }
      }
    },
    Bid: {
      teamByTeamId: {
        fragment: `... on Bid { teamId }`,
        resolve(bid, args, context, info) {
          return info.mergeInfo.delegateToSchema({
            schema: universeRemoteSchema,
            operation: 'query',
            fieldName: 'teamByTeamId',
            args: {
              teamId: bid.teamId,
            },
            context,
            info,
          })
        }
      }
    },
    Offer: {
      teamByBuyerTeamId: {
        fragment: `... on Offer { buyerTeamId }`,
        resolve(offer, args, context, info) {
          return info.mergeInfo.delegateToSchema({
            schema: universeRemoteSchema,
            operation: 'query',
            fieldName: 'teamByTeamId',
            args: {
              teamId: offer.buyerTeamId,
            },
            context,
            info,
          })
        }
      },
      playerByPlayerId: {
        fragment: `... on Offer { playerId }`,
        resolve(offer, args, context, info) {
          return info.mergeInfo.delegateToSchema({
            schema: universeRemoteSchema,
            operation: 'query',
            fieldName: 'playerByPlayerId',
            args: {
              playerId: offer.playerId,
            },
            context,
            info,
          })
        }
      }
    }
  };

  let schemas = [];
  universeRemoteSchema && schemas.push(universeRemoteSchema);
  marketRemoteSchema && schemas.push(marketRemoteSchema);
  relayRemoteSchema && schemas.push(relayRemoteSchema);
  notaryRemoteSchema && schemas.push(notaryRemoteSchema);
  schemas.push(linkTypeDefs);

  const schema = mergeSchemas({
    schemas,
    resolvers
  });

  const server = new ApolloServer({
    cors: enableCors ? true : false,
    schema,
  });

  server.listen().then(({ url }) => {
    console.log(`🚀  Server ready at ${url}`);
  });
};

const run = () => {
  main()
    .catch(e => {
      console.error(e);
      console.log("wainting ......");
      setTimeout(run, 3000);
    })
};

run();
