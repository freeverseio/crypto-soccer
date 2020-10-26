const { mailboxTypes } = require('../../config');
const GamelayerService = require('../GamelayerService');

const notifyNewHigherOffer = async ({
  destinatary,
  auctionId,
  text,
  metadata,
}) => {
  await GamelayerService.setMessage({
    destinatary,
    category: mailboxTypes.offer,
    auctionId,
    title: '',
    text,
    customImageUrl: '',
    metadata,
  });
};

module.exports = notifyNewHigherOffer;
