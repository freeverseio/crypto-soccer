const HorizonService = require('../services/HorizonService');

const allBidsResolver = async (parent, args, context, info, schema) => {
  const result = HorizonService.getAllBids();

  return result;
};

module.exports = allBidsResolver;
