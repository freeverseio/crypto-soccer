const GamelayerService = require('../GamelayerService');

const processAcceptedOffers = async ({ offerHistory }) => {
  await GamelayerService.setMessage({
    destinatary: offerHistory.buyer,
    category: 'offer',
    auctionId: offerHistory.auctionId,
    text: 'Blablbalba',
    customImageUrl: '',
    metadata: '',
  });
};

module.exports = processAcceptedOffers;
