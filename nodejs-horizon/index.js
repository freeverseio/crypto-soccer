const Web3 = require('web3');
const { GraphQLServer } = require('graphql-yoga')

const providerURL = 'ws://localhost:8545';

const web3 = new Web3(providerURL, null, {});
web3.eth.net.isListening().then(connected => {
  console.log('provider ' + providerURL + ' is connected: ' + connected);
})

const typeDefs = `
  type Query {
    hello(name: String): String!
  }
`

const resolvers = {
  Query: {
    hello: (_, { name }) => `Hello ${name || 'World'}`,
  },
}

const server = new GraphQLServer({ typeDefs, resolvers })

server.start(() => console.log('Server is running on localhost:4000'))
