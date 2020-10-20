const HorizonService = require('../services/HorizonService');

const allAuctionsResolver = async (parent, args, context, info, schema) => {
  const result = HorizonService.getAllAuctions();

  return result;
};

module.exports = allAuctionsResolver;
