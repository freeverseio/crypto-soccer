const { bidStates } = require('../../config');
const GamelayerService = require('../GamelayerService');
const HorizonService = require('../HorizonService');
const getTeamIdFromAuctionSeller = require('./getTeamIdFromAuctionSeller.js');

const processPayingAuction = async ({ auctionHistory }) => {
  const destinataryTeamId = await getTeamIdFromAuctionSeller({
    auction: auctionHistory,
  });

  if (destinataryTeamId) {
    await GamelayerService.setMessage({
      destinatary: destinataryTeamId,
      category: 'auction',
      auctionId: auctionHistory.id,
      text: 'auction_seller_sells',
      customImageUrl: '',
      metadata: '{}',
    });
  }

  const bids = await HorizonService.getBidsByAuctionId({
    auctionId: auctionHistory.id,
  });

  for (const bid of bids) {
    switch (bid.state) {
      case bidStates.paying:
        await GamelayerService.setMessage({
          destinatary: bid.teamId,
          category: 'auction',
          auctionId: auctionHistory.id,
          text: 'auction_buyer_wins_auction',
          customImageUrl: '',
          metadata: '{}',
        });
        break;
      default:
        await GamelayerService.setMessage({
          destinatary: bid.teamId,
          category: 'auction',
          auctionId: auctionHistory.id,
          text: 'auction_buyer_loses_auction',
          customImageUrl: '',
          metadata: '{}',
        });
        break;
    }
  }
};

module.exports = processPayingAuction;
