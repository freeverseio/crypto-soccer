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
  .option("-u, --universeUrl <url>", "graphql universe url", "http://localhost:4001/graphql")
  .option("-m, --marketUrl <url>", "graphql market url", "http://localhost:4002/graphql")
  .parse(process.argv);

const { universeUrl, marketUrl } = program;

console.log("--------------------------------------------------------");
console.log("universeUrl       : ", universeUrl);
console.log("marketUrl         : ", marketUrl);
console.log("--------------------------------------------------------");

const main = async () => {
  const universeRemoteSchema = await createRemoteSchema(universeUrl);
  const marketRemoteSchema = await createRemoteSchema(marketUrl);

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
      ): AuctionsConnection!
    }

    extend type Auction {
      playerByPlayerId: Player
    }

    extend type Bid {
      teamByTeamId: Team
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
  };

  const schema = mergeSchemas({
    schemas: [
      universeRemoteSchema,
      marketRemoteSchema,
      linkTypeDefs
    ],
    resolvers
  });

  const server = new ApolloServer({ schema });

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
