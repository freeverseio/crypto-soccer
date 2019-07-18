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
    createLeague: async (_, { nTeams, initBlock, step }) => {
      const gas = await leagues.methods.create(
        nTeams,
        initBlock,
        step
      ).estimateGas();
      await leagues.methods.create(
        nTeams,
        initBlock,
        step
      ).send({ from, gas });
      return true;
    },
  };
}

module.exports = Resolvers;