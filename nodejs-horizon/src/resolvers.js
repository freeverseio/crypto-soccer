const Web3 = require('web3');
const assetsContractJSON = require('../../truffle-core/build/contracts/Assets.json');

module.exports = function Resolvers(universe) {
  this.Query = {
    countTeams: () => universe.countTeams()
    // // teamById: async (_, params) => {
    //   const ids = await assetsContract.methods.getTeamPlayerIds(params.id).call();
    //   ids.forEach((part, index) => ids[index] = part.toString());
    //   return {
    //     id: params.id,
    //     name: await assetsContract.methods.getTeamName(params.id).call(),
    //     playerIds: ids
    //   }
    // },
    // allTeams: async () => {
    //   const count = await resolvers.Query.countTeams();
    //   let teams = [];
    //   for (let i = 1; i <= count; i++)
    //     teams.push(await resolvers.Query.teamById("", { id: i }));
    //   return teams;
    // }
  };

  this.Mutation = {
    createTeam: async (parent, args, context, info) => {
      await universe.createTeam(args.name, args.owner);
    }
  };

  this.Subscription = {
    // teamCreated: {
    //   subscribe: () => pubsub.asyncIterator([TEAM_CREATED])
    // }
  };
};