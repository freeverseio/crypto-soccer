const HorizonService = require('../HorizonService');
const GamelayerService = require('../GamelayerService');
const logger = require('../../logger');
const getTeamIdFromAuctionSeller = require('../auctions/getTeamIdFromAuctionSeller');
const { auctionStates } = require('../../config');

const processAcceptedBids = async ({ lastChecked }) => {
  const bids = await HorizonService.getLastAcceptedBidsHistories({
    lastChecked,
  });

  logger.info(`Processing Accepted Bids`);

  for (const bid of bids) {
    try {
      const auction = await HorizonService.getAuction({
        auctionId: bid.auctionId,
      });
      if (auction.state.toLowerCase() == auctionStates.assetFrozen) {
        const sellerTeamId = await getTeamIdFromAuctionSeller({ auction });
        const totalAmount = parseInt(auction.price) + parseInt(bid.extraPrice);
        const {
          name: bidderTeamName,
        } = await GamelayerService.getInfoFromTeamId({
          teamId: bid.teamId,
        });
        await GamelayerService.setMessage({
          destinatary: sellerTeamId,
          category: 'auction',
          auctionId: bid.auctionId,
          title: '',
          text: 'auction_seller_new_higher_bid',
          customImageUrl: '',
          metadata: `{"amount":"${totalAmount}", "playerId":"${auction.playerId}", "playerName":"${auction.playerByPlayerId.name}", "bidderTeamId":"${bid.teamId}", "bidderTeamName":"${bidderTeamName}"}`.replace(
            /"/g,
            '\\"'
          ),
        });

        for (const bidder of auction.bidsByAuctionId.nodes) {
          if (
            bidder.teamId != bid.teamId &&
            parseInt(bid.extraPrice) > parseInt(bidder.extraPrice)
          ) {
            await GamelayerService.setMessage({
              destinatary: bidder.teamId,
              category: 'auction',
              auctionId: bid.auctionId,
              title: '',
              text: 'auction_buyer_new_higher_bid',
              customImageUrl: '',
              metadata: `{"amount":"${totalAmount}", "playerId":"${auction.playerId}", "playerName":"${auction.playerByPlayerId.name}", "bidderTeamId":"${bid.teamId}", "bidderTeamName":"${bidderTeamName}"}`.replace(
                /"/g,
                '\\"'
              ),
            });
          }
        }
      }
    } catch (e) {
      logger.info(
        `Error processing accepted paying bid: ${JSON.stringify(bid)}`
      );
      logger.error(e);
    }
  }
};

module.exports = processAcceptedBids;
