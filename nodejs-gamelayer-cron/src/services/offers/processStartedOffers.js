const notifyNewHigherOffer = require('./notifyNewHigherOffer');
const GamelayerService = require('../GamelayerService');
const HorizonService = require('../HorizonService');

const processStartedOffers = async ({ offerHistory }) => {
  const offerers = await HorizonService.getOfferersByPlayerId({
    playerId: offerHistory.playerId,
  });

  switch (offerers.length) {
    case 1:
      await GamelayerService.setMessage({
        destinatary: offerHistory.teamId,
        category: 'offer',
        auctionId: offerHistory.auctionId,
        text: 'You received an offer for your player',
        customImageUrl: '',
        metadata: '',
      });
      break;
    default:
      await notifyNewHigherOffer({
        destinatary: offerHistory.seller,
        auctionId: offerHistory.auctionId,
      });
  }

  for (const offerer of offerers) {
    await notifyNewHigherOffer({
      destinatary: offerer.buyerTeamId,
      auctionId: offerHistory.auctionId,
    });
  }
};

module.exports = processStartedOffers;
