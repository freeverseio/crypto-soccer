const Web3 = require('web3');
const { ApolloServer, gql, PubSub } = require('apollo-server');
const typeDefs = require('./schema');
const Resolvers = require('./resolvers');

const assetsContractJSON = require('../../truffle-core/build/contracts/Assets.json');

const providerUrl = 'ws://localhost:8545';
const assetsContractAddress = '0xBaeb6C89EB37A467D8e54CCe11D1E093C5B18d6f';
const from = '0x9C33497cEc1E9603Ba65D3A8d5e59F543950d6Ef';
const gas = 6721975;

const web3 = new Web3(providerUrl, null, {});
const assetsContract = new web3.eth.Contract(assetsContractJSON.abi, assetsContractAddress);

const TEAM_CREATED = 'TEAM_CREATED';

const pubsub = new PubSub();
// const server = new GraphQLServer({ typeDefs, resolvers, context: { pubsub } });

assetsContract.events.TeamCreation()
  .on('data', (event) => {
    pubsub.publish(TEAM_CREATED, { teamCreated: event.returnValues.teamId.toString() });
  })
  .on('changed', (event) => {
    // remove event from local database
  })
  .on('error', console.error);

const resolvers = new Resolvers({provider: providerUrl, assetsContractAddress}).getResolvers();

// server.start(() => console.log('Server is running on localhost:4000'))

const server = new ApolloServer({ typeDefs, resolvers });

// This `listen` method launches a web-server.  Existing apps
// can utilize middleware options, which we'll discuss later.
server.listen().then(({ url }) => {
  console.log(`🚀  Server ready at ${url}`);
});