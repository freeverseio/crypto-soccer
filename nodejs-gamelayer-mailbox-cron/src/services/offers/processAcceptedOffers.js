const GamelayerService = require('../GamelayerService');
const HorizonService = require('../HorizonService');
const { mailboxTypes } = require('../../config');

const processAcceptedOffers = async ({ offerHistory }) => {
  const {
    teamId: playerTeamId,
    name: playerName,
  } = await HorizonService.getInfoFromPlayerId({
    playerId: offerHistory.playerId,
  });

  const { name: sellerTeamName } = await HorizonService.getInfoFromTeamId({
    teamId: playerTeamId,
  });

  await GamelayerService.setMessage({
    destinatary: offerHistory.buyerTeamId,
    category: mailboxTypes.offer,
    auctionId: offerHistory.auctionId,
    title: '',
    text: 'offer_buyer_offer_accepted',
    customImageUrl: '',
    metadata: `{"playerId":"${offerHistory.playerId}", "playerName": "${playerName}", "offerAmount": "${offerHistory.price}", "sellerTeamId": "${playerTeamId}", "sellerTeamName": "${sellerTeamName}"}`.replace(
      /"/g,
      '\\"'
    ),
  });
};

module.exports = processAcceptedOffers;
