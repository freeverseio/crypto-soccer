function Resolvers({
  market,
  from
}) {
  this.Query = {
    ping: () => true
  };

  this.Mutation = {
    transferFirstBotToAddr: async (_, { timezone, countryIdxInTimezone, address }) => {
      const gas = await market.methods
        .transferFirstBotToAddr(timezone, countryIdxInTimezone, address)
        .estimateGas();
      await market.methods
        .transferFirstBotToAddr(timezone, countryIdxInTimezone, address)
        .send({ from, gas });
      return true;
    }
  };
}

module.exports = Resolvers;