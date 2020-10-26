const { bidStates } = require('../../config');
const GamelayerService = require('../GamelayerService');
const HorizonService = require('../HorizonService');
const getTeamIdFromAuctionSeller = require('./getTeamIdFromAuctionSeller.js');

const processPayingAuction = async ({ auctionHistory }) => {
  const destinataryTeamId = await getTeamIdFromAuctionSeller({
    auction: auctionHistory,
  });

  const bids = await HorizonService.getBidsByAuctionId({
    auctionId: auctionHistory.id,
  });

  for (const bid of bids) {
    const { name: playerName } = await HorizonService.getInfoFromPlayerId({
      playerId: auctionHistory.playerId,
    });

    const { name: bidderTeamName } = await HorizonService.getInfoFromTeamId({
      teamId: bid.teamId,
    });

    const totalAmount =
      parseInt(auctionHistory.price) + parseInt(bid.extraPrice);

    switch (bid.state.toLowerCase()) {
      case bidStates.paying:
        await GamelayerService.setMessage({
          destinatary: bid.teamId,
          category: 'auction',
          auctionId: auctionHistory.id,
          text: 'auction_buyer_wins_auction',
          customImageUrl: '',
          metadata: `{"bidderTeamId":"${bid.teamId}", "bidderTeamName":"${bidderTeamName}", "amount": "${totalAmount}", "playerId": "${auctionHistory.playerId}", "playerName":"${playerName}", "paymentDeadline":"${bid.paymentDeadline}"}`.replace(
            /"/g,
            '\\"'
          ),
        });

        if (destinataryTeamId) {
          await GamelayerService.setMessage({
            destinatary: destinataryTeamId,
            category: 'auction',
            auctionId: auctionHistory.id,
            text: 'auction_seller_sells',
            customImageUrl: '',
            metadata: `{"bidderTeamId":"${bid.teamId}", "bidderTeamName":"${bidderTeamName}", "amount": "${totalAmount}", "playerId": "${auctionHistory.playerId}", "playerName":"${playerName}"}`.replace(
              /"/g,
              '\\"'
            ),
          });
        }
        break;
      case bidStates.accepted:
        const maxBid = bids.find(
          (b) =>
            b.state.toLowerCase() == bidStates.paid ||
            b.state.toLowerCase() == bidStates.paying
        );
        const {
          name: maxBidderTeamName,
        } = await HorizonService.getInfoFromTeamId({
          teamId: maxBid.teamId,
        });

        await GamelayerService.setMessage({
          destinatary: bid.teamId,
          category: 'auction',
          auctionId: auctionHistory.id,
          text: 'auction_buyer_loses_auction',
          customImageUrl: '',
          metadata: `{"amount": "${totalAmount}", "maxBidderTeamId":"${maxBid.teamId}", "maxBidderTeamName":"${maxBidderTeamName}", "playerId": "${auctionHistory.playerId}", "playerName":"${playerName}"}`.replace(
            /"/g,
            '\\"'
          ),
        });
        break;
    }
  }
};

module.exports = processPayingAuction;
