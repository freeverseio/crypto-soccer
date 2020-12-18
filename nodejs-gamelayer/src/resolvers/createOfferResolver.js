const { BidValidation } = require('../validations');

const computeAuctionId = ({ currencyId, price, rnd, playerId, validUntil, web3 }) => {
  const paramsSellerHiddenPrice = web3.eth.abi.encodeParameters(
    ['uint8', 'uint256', 'uint256'],
    [currencyId || 0, price || 0, rnd || 0]
  );
  const sellerHiddenPrice = web3.utils.soliditySha3(paramsSellerHiddenPrice);
  const params = web3.eth.abi.encodeParameters(
    ['bytes32', 'uint256', 'uint32'],
    [sellerHiddenPrice || '', playerId || 0, validUntil || 0]
  );
  return web3.utils.soliditySha3(params).slice(2);
};

const createOfferResolver = async (_, args, context, info, horizonRemoteSchema, web3) => {
  try {
    const {
      input: { buyerTeamId, rnd, signature, playerId, currencyId, price, validUntil },
    } = args;

    auctionId = await computeAuctionId({ currencyId, price, rnd, playerId, validUntil, web3 });

    const bidValidation = new BidValidation({ teamId: buyerTeamId, rnd: 0, auctionId, extraPrice: 0, signature, web3 });
    const isAllowed = await bidValidation.isAllowedToBid();

    if (!isAllowed) {
      return new Error('User not allowed to offer for that amount');
    } else {
      return info.mergeInfo.delegateToSchema({
        schema: horizonRemoteSchema,
        operation: 'mutation',
        fieldName: 'createOffer',
        args,
        context,
        info,
      });
    }
  } catch (e) {
    return e;
  }
};

module.exports = createOfferResolver;
