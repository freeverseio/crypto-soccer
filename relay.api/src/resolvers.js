const resolvers = (assets, from) => {
  return {
    Mutation: {
      transferFirstBotToAddr: async (_, { timezone, countryIdxInTimezone, address }) => {
        const gas = await assets.methods
          .transferFirstBotToAddr(timezone, countryIdxInTimezone, address)
          .estimateGas();
        await assets.methods
          .transferFirstBotToAddr(timezone, countryIdxInTimezone, address)
          .send({ from, gas });
        return true;
      },
    },
  };
};

module.exports = resolvers;