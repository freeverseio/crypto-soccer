const GamelayerService = require('../GamelayerService');

const processWithdrawableBySellerAuction = async ({ auctionHistory }) => {
  await GamelayerService.setMessage({
    destinatary: auctionHistory.seller,
    category: 'auction',
    auctionId: auctionHistory.id,
    text: 'Your auction is in state withdrawable by seller, go get paid.',
    customImageUrl: '',
    metadata: '{}',
  });
};

module.exports = processWithdrawableBySellerAuction;
