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
    createTeam: async (_, { name, owner }) => {
      const gas = await assets.methods.createTeam(name, owner).estimateGas();
      await assets.methods.createTeam(name, owner).send({ from: from, gas });
      return true;
    },
    createLeague: async (_, { initBlock, step, teamIds, tactics }) => {
      const count = await leagues.methods.leaguesCount().call();
      const id = count;
      const gas = await leagues.methods.create(
        id,
        initBlock,
        step,
        teamIds,
        tactics
      ).estimateGas();
      await leagues.methods.create(
        id,
        initBlock,
        step,
        teamIds,
        tactics
      ).send({ from, gas });
      return true;
    },
  };
}

module.exports = Resolvers;