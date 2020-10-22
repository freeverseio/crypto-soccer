const GamelayerService = require('../GamelayerService');
const { mailboxTypes } = require('../../config');

const processAcceptedOffers = async ({ offerHistory }) => {
  await GamelayerService.setMessage({
    destinatary: offerHistory.buyerTeamId,
    category: mailboxTypes.offer,
    auctionId: offerHistory.auctionId,
    text: 'offer_buyer_offer_accepted',
    customImageUrl: '',
    metadata: `{"playerId":"${offerHistory.playerId}"}`,
  });
};

module.exports = processAcceptedOffers;
