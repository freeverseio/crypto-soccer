const { ApolloServer, makeExecutableSchema, mergeSchemas } = require('apollo-server');
const { makeSchemaAndPlugin } = require("postgraphile-apollo-server");
const pg = require("pg");
const typeDefs = require('./schema');
const Resolvers = require('./resolvers');
const Web3 = require('web3');
const playerStateJSON = require('../../truffle-core/build/contracts/PlayerState.json');
const assetsJSON = require('../../truffle-core/build/contracts/Assets.json');
const leaguesJSON = require('../../truffle-core/build/contracts/Leagues.json');
const HDWalletProvider = require("truffle-hdwallet-provider");
const program = require('commander');
const version = require('../package.json').version;

// Parsing command line arguments
program
  .version(version)
  .option('-c, --config <path>', 'set config path. defaults to config.json')
  .parse(process.argv);

let configFile = "../";
if (typeof program.config !== 'undefined')
  configFile += program.config;
else
  configFile += "config.json";

console.log("Configuration file: " + configFile);
const config = require(configFile);
const {
  providerUrl,
  address,
  privateKey,
  statesContractAddress,
  assetsContractAddress,
  leaguesContractAddress
} = config;

console.log("--------------------------------------------------------");
console.log("providerUrl       : ", providerUrl);
console.log("account           : ", address);
console.log("🔥  account p.k.  : ", privateKey);
console.log("states address    : ", statesContractAddress);
console.log("assets address    : ", assetsContractAddress);
console.log("leagues address   : ", leaguesContractAddress);

console.log("--------------------------------------------------------");
const provider = new HDWalletProvider(privateKey, providerUrl);
const web3 = new Web3(provider, null, {});
const states = new web3.eth.Contract(playerStateJSON.abi, statesContractAddress);
const assets = new web3.eth.Contract(assetsJSON.abi, assetsContractAddress);
const leagues = new web3.eth.Contract(leaguesJSON.abi, leaguesContractAddress);

const pgPool = new pg.Pool({
  connectionString: 'postgres://freeverse:freeverse@localhost:5432/cryptosoccer'
});

makeSchemaAndPlugin(
  pgPool,
  "public", // PostgreSQL schema to use
  {
    retryOnInitFail: true,
    disableDefaultMutations: true,
    dynamicJson: true
  }
)
  .then(result => {
    const { schema, plugin } = result;

    const resolvers = new Resolvers({
      states,
      assets,
      leagues,
      from: address
    });

    const mutations = makeExecutableSchema({
      typeDefs: typeDefs,
      resolvers: resolvers
    });
    const mergedSchema = mergeSchemas({
      schemas: [schema, mutations]
    });

    console.log(mergedSchema);

    const server = new ApolloServer({
      schema: mergedSchema,
      // plugins: [plugin]
    });

    server.listen().then(({ url }) => {
      console.log(`🚀  Server ready at ${url}`);
    });
  })
  .catch(e => {
    console.error(e);
    process.exit(1);
  });

// const resolvers = new Resolvers({
//   states,
//   assets,
//   leagues,
//   from: address
// });

// // const pubsub = new PubSub();
// const server = new ApolloServer({ typeDefs, resolvers });

// // This `listen` method launches a web-server.  Existing apps
// // can utilize middleware options, which we'll discuss later.
// server.listen().then(({ url }) => {
//   console.log(`🚀  Server ready at ${url}`);
// });

/*
--------------------------------
Assets:         0xf60DAC49d2E0C7b3091A0423693757CEEeB642e5
States:         0xD5165DDd523F5dB1b20552fD949f149C363F417d
Engine:         0xe917715Db02C7355c06f2450042F2B25f5FEc77a
GameController: 0xC54CeBFeF6d3fed158C264f0a2dD6B46c89c0bbD
Leagues:        0xceA8d1CdB4518ca453039Cb4829518ff71DACE08
Stakers:        0x6c27FD6573DbCe335c6ee1480DFBC6FD4A0602b6
--------------------------------
*/


// const TEAM_CREATED = 'TEAM_CREATED';

// const server = new GraphQLServer({ typeDefs, resolvers, context: { pubsub } });

// assetsContract.events.TeamCreation()
//   .on('data', (event) => {
//     pubsub.publish(TEAM_CREATED, { teamCreated: event.returnValues.teamId.toString() });
//   })
//   .on('changed', (event) => {
//     // remove event from local database
//   })
//   .on('error', console.error);

