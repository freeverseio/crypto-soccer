const { insertMessages } = require('../repositories');
const HorizonService = require('../services/HorizonService.js');

const setBroadcastMessageResolver = async (
  _,
  { input: { category, auctionId, title, text, customImageUrl, metadata } }
) => {
  try {
    const teamIdsArray = await HorizonService.getAllTeamIds();
    const messages = teamIdsArray.map((teamIdObj) => {
      const { teamId } = teamIdObj;
      return { destinatary: teamId, category, auctionId, title, text, customImageUrl, metadata };
    });
    if (!messages.length) {
      return false;
    }
    await insertMessages(messages);
    return true;
  } catch (e) {
    return e;
  }
};

module.exports = setBroadcastMessageResolver;
