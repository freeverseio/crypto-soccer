function Resolvers({
  states,
  assets,
  leagues,
  from
}) {
  this.Query = {
    ping: () => true
  };

  this.Mutation = {
    createLeague: async (_, { nTeams, initBlock, step }) => {
      const gas = await leagues.methods
        .create(nTeams, initBlock, step)
        .estimateGas();
      await leagues.methods
        .create(nTeams, initBlock, step)
        .send({ from, gas });
      return true;
    },
    transferTeam: async (_, { teamId, owner }) => {
      const gas = await assets.methods.transferTeam(teamId, owner).estimateGas();
      await assets.methods.transferTeam(teamId, owner).send({ from, gas });
      return true;
    }
  };
}

module.exports = Resolvers;