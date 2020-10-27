const GamelayerService = require('../GamelayerService');
const HorizonService = require('../HorizonService');

const processRejectedOffers = async ({ offerHistory }) => {
  const {
    teamId: playerTeamId,
    name: playerName,
  } = await HorizonService.getInfoFromPlayerId({
    playerId: offerHistory.playerId,
  });

  const { name: sellerTeamName } = await GamelayerService.getInfoFromTeamId({
    teamId: playerTeamId,
  });

  await GamelayerService.setMessage({
    destinatary: offerHistory.buyerTeamId,
    category: 'offer',
    auctionId: offerHistory.auctionId,
    title: '',
    text: 'offer_buyer_offer_rejected',
    customImageUrl: '',
    metadata: `{"playerId":"${offerHistory.playerId}", "playerName":"${playerName}", "sellerTeamId":"${playerTeamId}", "sellerTeamName":"${sellerTeamName}"}`.replace(
      /"/g,
      '\\"'
    ),
  });
};

module.exports = processRejectedOffers;
