const { mailboxTypes } = require('../../config');
const GamelayerService = require('../GamelayerService');

const notifyNewHigherOffer = async ({ destinatary, auctionId }) => {
  await GamelayerService.setMessage({
    destinatary,
    category: mailboxTypes.offer,
    auctionId,
    text: 'Blablbalba',
    customImageUrl: '',
    metadata: '',
  });
};

module.exports = notifyNewHigherOffer;
