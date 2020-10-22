const GamelayerService = require('../GamelayerService');

const processRejectedOffers = async ({ offerHistory }) => {
  await GamelayerService.setMessage({
    destinatary: offerHistory.buyerTeamId,
    category: 'offer',
    auctionId: offerHistory.auctionId,
    text: 'offer_buyer_offer_rejected',
    customImageUrl: '',
    metadata: `{"playerId":"${offerHistory.playerId}"}`,
  });
};

module.exports = processRejectedOffers;
