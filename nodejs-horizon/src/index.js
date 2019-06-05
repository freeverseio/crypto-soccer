const { ApolloServer, PubSub } = require('apollo-server');
const typeDefs = require('./schema');
const Resolvers = require('./resolvers');
const Web3 = require('web3');
const playerStateJSON = require('../../truffle-core/build/contracts/PlayerState.json');
const assetsJSON = require('../../truffle-core/build/contracts/Assets.json');
const leaguesJSON = require('../../truffle-core/build/contracts/Leagues.json');
const HDWalletProvider = require("truffle-hdwallet-provider");

const address = '0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7';
const privateKey = "3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54";
const providerUrl = "https://devnet.busyverse.com/web3";
const provider = new HDWalletProvider("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54", "https://devnet.busyverse.com/web3");

const web3 = new Web3(provider, null, {});
const states = new web3.eth.Contract(playerStateJSON.abi, '0x0836cB83c11Ce40C77eCF77a541a32c26C146b79');
const assets = new web3.eth.Contract(assetsJSON.abi, '0xd2A1716E791C73C6E8CC589E75452C85B64bed7a');
const leagues = new web3.eth.Contract(leaguesJSON.abi, '0x3629436292F52830dBd467A95eAdd8BB3d84f9eb');

const resolvers = new Resolvers({
  states,
  assets,
  leagues,
  from: address
});

const pubsub = new PubSub();
const server = new ApolloServer({ typeDefs, resolvers });

// This `listen` method launches a web-server.  Existing apps
// can utilize middleware options, which we'll discuss later.
server.listen().then(({ url }) => {
  console.log(`🚀  Server ready at ${url}`);
});

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

