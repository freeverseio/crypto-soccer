const { BidValidation } = require('../validations');
const CustomError = require('./CustomError');

const createBidResolver = async (_, args, context, info, horizonRemoteSchema, web3) => {
  try {
    const {
      input: { teamId, rnd, auctionId, extraPrice, signature },
    } = args;
    const bidValidation = new BidValidation({ teamId, rnd, auctionId, extraPrice, signature, web3 });
    const isAllowed = await bidValidation.isAllowedToBid();

    if (!isAllowed) {
      err = new CustomError('100', 'User not allowed to bid for that amount');
      return err;
    } else {
      return info.mergeInfo.delegateToSchema({
        schema: horizonRemoteSchema,
        operation: 'mutation',
        fieldName: 'createBid',
        args,
        context,
        info,
      });
    }
  } catch (e) {
    return e;
  }
};

module.exports = createBidResolver;
