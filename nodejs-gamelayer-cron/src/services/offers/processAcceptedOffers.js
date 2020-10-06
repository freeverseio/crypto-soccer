const GamelayerService = require('../GamelayerService');
const { mailboxTypes } = require('../../config');

const processAcceptedOffers = async ({ offerHistory }) => {
  await GamelayerService.setMessage({
    destinatary: offerHistory.buyerTeamId,
    category: mailboxTypes.offer,
    auctionId: offerHistory.auctionId,
    text: 'Blablbalba',
    customImageUrl: '',
    metadata: '',
  });
};

module.exports = processAcceptedOffers;
