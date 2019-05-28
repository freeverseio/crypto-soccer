const Web3 = require('web3');
const { GraphQLServer, PubSub } = require('graphql-yoga');

const assetsContractJSON = require('../../truffle-core/build/contracts/Assets.json');

const providerUrl = 'ws://localhost:8545';
const assetsContractAddress = '0xBaeb6C89EB37A467D8e54CCe11D1E093C5B18d6f';
const from = '0x9C33497cEc1E9603Ba65D3A8d5e59F543950d6Ef';
const gas = 6721975;

const web3 = new Web3(providerUrl, null, {});
const assetsContract = new web3.eth.Contract(assetsContractJSON.abi, assetsContractAddress);

const TEAM_CREATED = 'TEAM_CREATED';

const typeDefs = `
  type Query {
    settings: Settings!
    countTeams: String!
    teamById(id: ID!): Team
    allTeams: [Team]
  }

  type Mutation {
    createTeam(name: String!, owner: String!): String
  }

  type Subscription {
    teamCreated: ID!
  }

  type Settings {
    providerUrl: String
    assetsContractAddress: String
    from: String
    gas: String
  }

  type Team {
    id: ID!
    name: String!
    playerIds: [ID!]
  }
`;

const resolvers = {
  Query: {
    settings: () => ({
      providerUrl: web3.currentProvider.connection._url,
      assetsContractAddress: assetsContractAddress,
      from,
      gas
    }),
    countTeams: async () => {
      const count = await assetsContract.methods.countTeams().call();
      return count.toString();
    },
    teamById: async (_, params) => {
      const ids = await assetsContract.methods.getTeamPlayerIds(params.id).call();
      ids.forEach((part, index) => ids[index] = part.toString());
      return {
        id: params.id,
        name: await assetsContract.methods.getTeamName(params.id).call(),
        playerIds: ids
      }
    },
    allTeams: async () => {
      const count = await resolvers.Query.countTeams();
      let teams = [];
      for (let i=1 ; i <= count ; i++)
        teams.push(await resolvers.Query.teamById("", {id: i}));
      return teams;
    }
  },
  Mutation: {
    createTeam: (_, params) => {
      assetsContract.methods.createTeam(params.name, params.owner).send({ from, gas });
    }
  },
  Subscription: {
    teamCreated: {
      subscribe: () => pubsub.asyncIterator([TEAM_CREATED])
    }
  },
}

const pubsub = new PubSub();
const server = new GraphQLServer({ typeDefs, resolvers, context: { pubsub } });

assetsContract.events.TeamCreation()
  .on('data', (event) => {
    pubsub.publish(TEAM_CREATED, { teamCreated: event.returnValues.teamId.toString() });
  })
  .on('changed', (event) => {
    // remove event from local database
  })
  .on('error', console.error);

server.start(() => console.log('Server is running on localhost:4000'))
