const GamelayerService = require('../Gamelayerservice');

const processRejectedOffers = async ({ offerHistory }) => {
  await GamelayerService.setMessage({
    destinatary: offerHistory.buyer,
    category: 'offer',
    auctionId: offerHistory.auctionId,
    text: 'Blablbalba',
    customImageUrl: '',
    metadata: '',
  });
};

module.exports = processRejectedOffers;
