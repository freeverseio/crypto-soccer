const Web3 = require('web3');
const {
  GraphQLServer
} = require('graphql-yoga');

const assetsContractJSON = require('../truffle-core/build/contracts/Assets.json');

const providerUrl = 'ws://localhost:8545';
const assetsContractAddress = '0x5841d35b6580b1f7599b2e1157CAd16690A86f16';
const from = '0x9C33497cEc1E9603Ba65D3A8d5e59F543950d6Ef';
const gas = 6721975;

const web3 = new Web3(providerUrl, null, {});
const assetsContract = new web3.eth.Contract(assetsContractJSON.abi, assetsContractAddress);

const typeDefs = `
  type Query {
    getSettings: Settings!
    getTeamCount: String!
    getTeam(id: ID!): Team
  }

  type Mutation {
    createTeam(name: String!, owner: String!): String
  }

  type Settings {
    providerUrl: String
    assetsContractAddress: String
    from: String
    gas: String
  }

  type Provider {
    url: String
    isListening: Boolean!
  }

  type Team {
    id: ID!
    name: String!
  }
`
const resolvers = {
  Query: {
    getSettings: () => ({
      providerUrl: web3.currentProvider.connection._url,
      assetsContractAddress: assetsContractAddress,
      from,
      gas
    }),
    getTeamCount: async () => {
      const count = await assetsContract.methods.countTeams().call()
      return count.toString();
    },
    getTeam: async (_, params) => ({
      id: params.id,
      name: await assetsContract.methods.getTeamName(params.id).call()
    })
  },
  Mutation: {
    createTeam: (_, params) => {
      assetsContract.methods.createTeam(params.name, params.owner).send({
        from,
        gas
      });
    }
  }
}

const server = new GraphQLServer({
  typeDefs,
  resolvers
})

server.start(() => console.log('Server is running on localhost:4000'))
