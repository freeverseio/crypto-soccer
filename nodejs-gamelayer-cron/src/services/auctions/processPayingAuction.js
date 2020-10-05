const { bidStates } = require('../../config');
const GamelayerService = require('../Gamelayerservice');
const HorizonService = require('../HorizonService');

const processPayingAuction = async ({ auctionHistory }) => {
  await GamelayerService.setMessage({
    destinatary: auctionHistory.seller,
    category: 'auction',
    auctionId: auctionHistory.id,
    text: 'The auction for this player has ended. You sold him!',
    customImageUrl: '',
    metadata: '',
  });

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
          text: 'You won the auction, remember to pay',
          customImageUrl: '',
          metadata: '',
        });
        break;
      default:
        await GamelayerService.setMessage({
          destinatary: bid.teamId,
          category: 'auction',
          auctionId: auctionHistory.id,
          text: 'You lost the auction',
          customImageUrl: '',
          metadata: '',
        });
        break;
    }
  }
};

module.exports = processPayingAuction;
