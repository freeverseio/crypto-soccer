const Web3 = require('web3');
const assetsContractJSON = require('../../truffle-core/build/contracts/Assets.json');
const assetsContractAddress = '0xBaeb6C89EB37A467D8e54CCe11D1E093C5B18d6f';

module.exports = class Resolvers {
  constructor({
    provider,
    assetsContractAddress
  }) {
    this.web3 = new Web3(provider, null, {});
    this.assetsContract = new this.web3.eth.Contract(assetsContractJSON.abi, assetsContractAddress);
  }

  getResolvers() {
    const { web3, assetsContract } = this;
    return {
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
          for (let i = 1; i <= count; i++)
            teams.push(await resolvers.Query.teamById("", { id: i }));
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
  }
};