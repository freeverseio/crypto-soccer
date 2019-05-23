const Web3 = require('web3');
const {
  GraphQLServer
} = require('graphql-yoga');
const assetsContractJSON = require('../truffle-core/build/contracts/Assets.json');
const assetsContractAddress = '0x14ebD51cD831f8B371b19ad571FaCDa3655004a4';
const providerURL = 'ws://localhost:8545';
const web3 = new Web3(providerURL, null, {});

web3.eth.net.isListening()
  .then(connected => {
    console.log('provider ' + providerURL + ' is connected: ' + connected);
    const assetsContract = new web3.eth.Contract(assetsContractJSON.abi, assetsContractAddress);

    const typeDefs = `
        type Query {
          getProvider: Provider!
          getAssetsContractAddress: String
          getTeamCount: String!
        }

        type Provider {
          url: String
          isListening: Boolean!
        }
      `
    const resolvers = {
      Query: {
        getProvider: async () => ({
          url: web3.currentProvider.connection._url,
          isListening: await web3.eth.net.isListening()
        }),
        getAssetsContractAddress: () => (assetsContract.address),
        getTeamCount: async () => {
          const count = await assetsContract.methods.countTeams().call()
          return count.toString();
        }
      },
    }

    const server = new GraphQLServer({
      typeDefs,
      resolvers
    })

    server.start(() => console.log('Server is running on localhost:4000'))
  })
  .catch(console.log);