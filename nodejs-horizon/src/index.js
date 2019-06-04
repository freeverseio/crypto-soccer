const { ApolloServer, PubSub } = require('apollo-server');
const Universe = require('./universe/Universe');
const typeDefs = require('./schema');
const Resolvers = require('./resolvers');

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

const providerUrl = 'ws://localhost:8545';
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
  playerStateAddress: '0xD5165DDd523F5dB1b20552fD949f149C363F417d',
  assetsAddress: '0xf60DAC49d2E0C7b3091A0423693757CEEeB642e5',
  leaguesAddress: '0xceA8d1CdB4518ca453039Cb4829518ff71DACE08',
  from
});
const resolvers = new Resolvers({
  provider: providerUrl,
  playerStateAddress: '0xD5165DDd523F5dB1b20552fD949f149C363F417d',
  assetsAddress: '0xf60DAC49d2E0C7b3091A0423693757CEEeB642e5',
  leaguesAddress: '0xceA8d1CdB4518ca453039Cb4829518ff71DACE08',
  from
});

const server = new ApolloServer({ typeDefs, resolvers });

// This `listen` method launches a web-server.  Existing apps
// can utilize middleware options, which we'll discuss later.
server.listen().then(({ url }) => {
  console.log(`🚀  Server ready at ${url}`);
});