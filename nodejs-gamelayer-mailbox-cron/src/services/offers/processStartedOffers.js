const notifyNewHigherOffer = require('./notifyNewHigherOffer');
const GamelayerService = require('../GamelayerService');
const HorizonService = require('../HorizonService');
const { mailboxTypes } = require('../../config');

const processStartedOffers = async ({ offerHistory }) => {
  const offerers = await HorizonService.getOfferersByPlayerId({
    playerId: offerHistory.playerId,
  });

  const {
    teamId: playerTeamId,
    name,
  } = await HorizonService.getInfoFromPlayerId({
    playerId: offerHistory.playerId,
  });

  const { name: offererTeamName } = await HorizonService.getInfoFromTeamId({
    teamId: offerHistory.buyerTeamId,
  });

  switch (offerers.length) {
    case 1:
      await GamelayerService.setMessage({
        destinatary: playerTeamId,
        category: mailboxTypes.offer,
        auctionId: offerHistory.auctionId,
        text: 'offer_seller_offer_received',
        customImageUrl: '',
        metadata: `{"playerId":"${offerHistory.playerId}", "playerName": "${name}", "offerAmount": "${offerHistory.price}", "offererTeamId": "${offerHistory.buyerTeamId}", "offererTeamName": "${offererTeamName}"}`.replace(
          /"/g,
          '\\"'
        ),
      });
      break;
    default:
      await notifyNewHigherOffer({
        destinatary: playerTeamId,
        auctionId: offerHistory.auctionId,
        text: 'offer_seller_higher_offer_received',
        metadata: `{"playerId":"${offerHistory.playerId}", "playerName": "${name}", "offerAmount": "${offerHistory.price}", "offererTeamId": "${offerHistory.buyerTeamId}", "offererTeamName": "${offererTeamName}"}`.replace(
          /"/g,
          '\\"'
        ),
      });
  }

  for (const offerer of offerers) {
    await notifyNewHigherOffer({
      destinatary: offerer.buyerTeamId,
      auctionId: offerHistory.auctionId,
      text: 'offer_buyer_higher_offer',
      metadata: `{"playerId":"${offerHistory.playerId}", "playerName": "${name}", "offerAmount": "${offerHistory.price}", "offererTeamId": "${offerHistory.buyerTeamId}", "offererTeamName": "${offererTeamName}"}`.replace(
        /"/g,
        '\\"'
      ),
    });
  }
};

module.exports = processStartedOffers;
