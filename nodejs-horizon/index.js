const Web3 = require('web3');
const { GraphQLServer } = require('graphql-yoga')

const providerURL = 'ws://localhost:8545';

const web3 = new Web3(providerURL, null, {});
web3.eth.net.isListening().then(connected => {
  console.log('provider ' + providerURL + ' is connected: ' + connected);
})

const typeDefs = `
  type Query {
    hello(name : String): String!
    team: Team!
    provider: Provider!
  }

  type Provider {
    url: String!
    isListening: Boolean!
  }

  type Team {
    name: String!
  }
`

const resolvers = {
  Query: {
    hello: (_, { name }) => `Hello ${name || 'World'}`,
    team: () =>  (''),
    provider: () => ('')
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
