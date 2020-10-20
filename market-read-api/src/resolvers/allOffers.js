const HorizonService = require('../services/HorizonService');

const allOffersResolver = async (parent, args, context, info, schema) => {
  const result = HorizonService.getAllOffers();

  return result;
};

module.exports = allOffersResolver;
