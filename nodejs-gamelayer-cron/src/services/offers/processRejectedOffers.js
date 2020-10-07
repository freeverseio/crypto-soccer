const GamelayerService = require('../GamelayerService');

const processRejectedOffers = async ({ offerHistory }) => {
  await GamelayerService.setMessage({
    destinatary: offerHistory.buyerTeamId,
    category: 'offer',
    auctionId: offerHistory.auctionId,
    text: 'Your offer for this player has been rejected',
    customImageUrl: '',
    metadata: '{}',
  });
};

module.exports = processRejectedOffers;
