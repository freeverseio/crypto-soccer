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
      saleOrderByPlayerId: PlayerSaleOrder
    }
  `;

  const resolvers = {
    Player: {
      saleOrderByPlayerId: {
        fragment: `... on Player { id }`,
        resolve(player, args, context, info) {
          return info.mergeInfo.delegateToSchema({
            schema: marketRemoteSchema,
            operation: 'query',
            fieldName: 'playerSaleOrderByPlayerid',
            args: {
              playerid: player.id,
            },
            context,
            info,
          })
        }
      }
    }
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

try {
  main();
} catch (e) {
  console.log(e, e.message, e.stack);
}
