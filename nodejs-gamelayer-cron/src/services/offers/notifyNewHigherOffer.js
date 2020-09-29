const notifyNewHigherOffer = async ({ destinatary, auctionId }) => {
  await GamelayerService.setMessage({
    destinatary,
    category: 'offer',
    auctionId,
    text: 'Blablbalba',
    customImageUrl: '',
    metadata: '',
  });
};

module.exports = notifyNewHigherOffer;
