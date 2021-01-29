const { OfferValidation } = require('../validations');
const CustomError = require('./CustomError');

const createOfferResolver = async (_, args, context, info, horizonRemoteSchema, web3) => {
  try {
    const {
      input: { buyerTeamId, rnd, signature, playerId, currencyId, price, validUntil },
    } = args;

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
      err = new CustomError('101', 'User not allowed to offer for that amount');
      return err;
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
