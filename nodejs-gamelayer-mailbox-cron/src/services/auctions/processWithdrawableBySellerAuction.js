const GamelayerService = require('../GamelayerService');
const HorizonService = require('../HorizonService');
const getTeamIdFromAuctionSeller = require('./getTeamIdFromAuctionSeller.js');

const processWithdrawableBySellerAuction = async ({ auctionHistory }) => {
  const destinataryTeamId = await getTeamIdFromAuctionSeller({
    auction: auctionHistory,
  });

  if (destinataryTeamId) {
    const paidBid = await HorizonService.getPaidBidByAuctionId({
      auctionId: auctionHistory.id,
    });
    const { name: playerName } = await HorizonService.getInfoFromPlayerId({
      playerId: auctionHistory.playerId,
    });

    const { name: bidderTeamName } = await HorizonService.getInfoFromTeamId({
      teamId: paidBid.teamId,
    });

    const totalAmount =
      parseInt(auctionHistory.price) + parseInt(paidBid.extraPrice);

    await GamelayerService.setMessage({
      destinatary: destinataryTeamId,
      category: 'auction',
      auctionId: auctionHistory.id,
      text: 'auction_seller_gets_paid',
      customImageUrl: '',
      metadata: `{"bidderTeamId":"${paidBid.teamId}", "bidderTeamName":"${bidderTeamName}", "amount": "${totalAmount}", "playerId": "${auctionHistory.playerId}", "playerName":"${playerName}"}`,
    });
  }
};

module.exports = processWithdrawableBySellerAuction;
