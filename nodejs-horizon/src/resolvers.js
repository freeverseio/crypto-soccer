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