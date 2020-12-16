const { upsertOwnerMaxBidAllowed } = require('../repositories');
const { TOKEN } = require('../config');

const setMaxBidAllowedByOwnerResolver = async (_, { input: { owner, maxBidAllowed, token } }) => {
  try {
    if (token === TOKEN) {
      return Error('Token not allowed');
    }
    await upsertOwnerMaxBidAllowed({ owner, maxBidAllowed });

    return true;
  } catch (e) {
    return e;
  }
};

module.exports = setMaxBidAllowedByOwnerResolver;
