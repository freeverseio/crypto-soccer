const Web3 = require('web3');
const { GraphQLServer } = require('graphql-yoga')

const providerURL = 'ws://localhost:8545';
const web3 = new Web3(providerURL, null, {});

web3.eth.net.isListening()
  .then(connected => {
    console.log('provider ' + providerURL + ' is connected: ' + connected);
  })
  .catch(console.log);

const typeDefs = `
  type Query {
    hello(name : String): String!
    team: Team!
    getProvider: Provider!
    getAssetsContractAddress: String
  }

  type Provider {
    url: String
    isListening: Boolean!
  }

  type Mutation {
    setProviderUrl(url: String!): Provider!
  }

  type Team {
    name: String!
  }
`

const resolvers = {
  Query: {
    hello: (_, { name }) => `Hello ${name || 'World'}`,
    team: (_) =>  (''),
    getProvider: async () => ({
      url: web3.currentProvider.connection._url,
      isListening: await web3.eth.net.isListening()
    }),
    getAssetsContractAddress: () => (assetsContractAddress)
  },
  Mutation: {
    setProviderUrl: () => (console.log("setProviderUrl"))
  },
  Provider: {
    url: () => (web3.currentProvider.connection._url),
    isListening: async () => (await web3.eth.net.isListening()),
  },
  Team: {
    name: () => 'Ciao qui'
  }
}

const server = new GraphQLServer({ typeDefs, resolvers })

server.start(() => console.log('Server is running on localhost:4000'))
