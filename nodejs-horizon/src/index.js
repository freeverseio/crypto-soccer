const { ApolloServer, PubSub } = require('apollo-server');
const Universe = require('./universe/Universe');
const typeDefs = require('./schema');
const Resolvers = require('./resolvers');

const providerUrl = 'ws://localhost:8545';
const assetsContractAddress = '0xBaeb6C89EB37A467D8e54CCe11D1E093C5B18d6f';
const from = '0x9C33497cEc1E9603Ba65D3A8d5e59F543950d6Ef';

const TEAM_CREATED = 'TEAM_CREATED';

const pubsub = new PubSub();
// const server = new GraphQLServer({ typeDefs, resolvers, context: { pubsub } });

// assetsContract.events.TeamCreation()
//   .on('data', (event) => {
//     pubsub.publish(TEAM_CREATED, { teamCreated: event.returnValues.teamId.toString() });
//   })
//   .on('changed', (event) => {
//     // remove event from local database
//   })
//   .on('error', console.error);

const universe = new Universe({
  provider: providerUrl,
  assetsAddress: assetsContractAddress,
  playerStateAddress: '0x6367bd320b0Ce4C99f58685E6ca61F0F9660cdb9',
  from
});
const resolvers = new Resolvers(universe);

const server = new ApolloServer({ typeDefs, resolvers });

// This `listen` method launches a web-server.  Existing apps
// can utilize middleware options, which we'll discuss later.
server.listen().then(({ url }) => {
  console.log(`🚀  Server ready at ${url}`);
});