const GamelayerService = require('../GamelayerService');

const processRejectedOffers = async ({ offerHistory }) => {
  await GamelayerService.setMessage({
    destinatary: offerHistory.buyerTeamId,
    category: 'offer',
    auctionId: offerHistory.auctionId,
    text: 'Blablbalba',
    customImageUrl: '',
    metadata: '',
  });
};

module.exports = processRejectedOffers;
