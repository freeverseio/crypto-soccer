const { OfferValidation } = require('../validations');

const createOfferResolver = async (_, args, context, info, horizonRemoteSchema, web3) => {
  try {
    const {
      input: { buyerTeamId, rnd, signature, playerId, currencyId, price, validUntil },
    } = args;

    auctionId = await computeAuctionId({ currencyId, price, rnd, playerId, validUntil, web3 });

    const offerValidation = new OfferValidation({
      currencyId,
      playerId,
      price,
      validUntil,
      buyerTeamId,
      rnd,
      signature,
      web3,
    });
    const isAllowed = await offerValidation.isAllowedToOffer(true);

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
