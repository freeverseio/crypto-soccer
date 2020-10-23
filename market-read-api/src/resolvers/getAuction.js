const HorizonService = require('../services/HorizonService');

const getAuctionResolver = async (parent, args, context, info, schema) => {
  const { id } = args;
  const result = HorizonService.getAuction({ id });

  return result;
};

module.exports = getAuctionResolver;
