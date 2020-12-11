const upsertTeamGetSocialId = require('../repositories/upsertTeamGetSocialId');
const { GetSocialIdValidation } = require('../validations');

const setGetSocialIdResolver = async (_, args, context, info, web3) => {
  try {
    const {
      input: { teamId, getSocialId, signature },
    } = args;
    const getSocialIdValidation = new GetSocialIdValidation({ teamId, getSocialId, signature, web3 });
    const isSignerOwner = await getSocialIdValidation.isSignerOwner();

    if (!isSignerOwner) {
      return new Error('Signer is not the owner of the team');
    } else {
      await upsertTeamGetSocialId({ teamId, getSocialId });
      return true;
    }
  } catch (e) {
    return e;
  }
};

module.exports = setGetSocialIdResolver;
