const { insertMessage } = require('../repositories');
const HorizonService = require('../services/HorizonService.js');

const setBroadcastMessageResolver = async (
  _,
  { input: { category, auctionId, text, customImageUrl, metadata } }
) => {
  try {
    const teamIds = await HorizonService.getAllTeamIds();
    for (const teamId of teamIds) {
      const idFromDb = await insertMessage({
        destinatary: teamId,
        category,
        auctionId,
        text,
        customImageUrl,
        metadata,
      });
    }

    return true;
  } catch (e) {
    return e;
  }
};

module.exports = setBroadcastMessageResolver;
