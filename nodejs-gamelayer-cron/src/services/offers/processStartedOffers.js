const notifyNewHigherOffer = require('./notifyNewHigherOffer');

const processStartedOffers = async ({ offerHistory }) => {
  await notifyNewHigherOffer({
    destinatary: offerHistory.seller,
    auctionId: offerHistory.auctionId,
  });
  const offerers = await HorizonService.getOfferersByPlayerId({
    playerId: offerHistory.playerId,
  });
  for (const offerer of offerers) {
    await notifyNewHigherOffer({
      destinatary: offerer,
      auctionId: offerHistory.auctionId,
    });
  }
};

module.exports = processStartedOffers;
