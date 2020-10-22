const GamelayerService = require('../GamelayerService');
const getTeamIdFromAuctionSeller = require('./getTeamIdFromAuctionSeller.js');

const processWithdrawableBySellerAuction = async ({ auctionHistory }) => {
  const destinataryTeamId = await getTeamIdFromAuctionSeller({
    auction: auctionHistory,
  });

  if (destinataryTeamId) {
    await GamelayerService.setMessage({
      destinatary: destinataryTeamId,
      category: 'auction',
      auctionId: auctionHistory.id,
      text: 'auction_seller_gets_paid',
      customImageUrl: '',
      metadata: '{}',
    });
  }
};

module.exports = processWithdrawableBySellerAuction;
