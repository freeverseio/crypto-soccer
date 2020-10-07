const { mailboxTypes } = require('../../config');
const GamelayerService = require('../GamelayerService');

const notifyNewHigherOffer = async ({ destinatary, auctionId }) => {
  await GamelayerService.setMessage({
    destinatary,
    category: mailboxTypes.offer,
    auctionId,
    text: 'Your player has received a new higher offer!',
    customImageUrl: '',
    metadata: '',
  });
};

module.exports = notifyNewHigherOffer;
