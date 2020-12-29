const dayjs = require('dayjs');
const { bidStates } = require('../../config');
const GamelayerService = require('../GamelayerService');
const HorizonService = require('../HorizonService');
const getTeamIdFromAuctionSeller = require('./getTeamIdFromAuctionSeller.js');
const logger = require('../../logger');

const getFormattedPaymentDeadline = async ({ paymentDeadline }) => {
  logger.debug(
    `unix now: ${dayjs().unix()} || bids[0].paymentDeadline: ${
      bids[0].paymentDeadline
    } || deadline: ${deadline}`
  );
  let deadline = dayjs().add(46, 'hour').unix();

  if (bids[0].paymentDeadline && bids[0].paymentDeadline > dayjs().unix()) {
    deadline = bids[0].paymentDeadline;
  }

  return deadline;
};

const processPayingAuction = async ({ auctionHistory }) => {
  const destinataryTeamId = await getTeamIdFromAuctionSeller({
    auction: auctionHistory,
  });

  const bids = await HorizonService.getBidsByAuctionId({
    auctionId: auctionHistory.id,
  });
  logger.debug(
    `ProcessPayingAuction: Auction history: ${JSON.stringify(
      auctionHistory
    )}\n\nBids: ${JSON.stringify(bids)}`
  );

  let maxBid =
    bids.length == 1
      ? bids[0]
      : bids.find(
          (b) =>
            b.state.toLowerCase() == bidStates.paid ||
            b.state.toLowerCase() == bidStates.paying
        );

  if (!maxBid) {
    const maxExtraPrice = Math.max(...bids.map((b) => b.extraPrice));
    maxBid = bids.find((b) => b.extraPrice == maxExtraPrice);
    const maxBidIndex = bids.findIndex((b) => b.extraPrice == maxExtraPrice);
    logger.debug(`Maxbid(bids[${maxBidIndex}]): ${JSON.stringify(maxBid)}`);
    bids[maxBidIndex].state = 'PAYING';
    bids[maxBidIndex].paymentDeadline = getFormattedPaymentDeadline({
      paymentDeadline: bids[maxBidIndex].paymentDeadline,
    });
  }

  if (bids.length == 1) {
    bids[0].state = 'PAYING';
    bids[0].paymentDeadline = getFormattedPaymentDeadline({
      paymentDeadline: bids[0].paymentDeadline,
    });
  }

  const { name: maxBidderTeamName } = await GamelayerService.getInfoFromTeamId({
    teamId: maxBid.teamId,
  });
  const buyerLosesAuctionTeamsNotified = [];

  for (const bid of bids) {
    const { name: playerName } = await HorizonService.getInfoFromPlayerId({
      playerId: auctionHistory.playerId,
    });

    const { name: bidderTeamName } = await GamelayerService.getInfoFromTeamId({
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
          title: '',
          text: 'auction_buyer_wins_auction',
          customImageUrl: '',
          metadata: `{"bidderTeamId":"${bid.teamId}", "bidderTeamName":"${bidderTeamName}", "amount":"${totalAmount}", "playerId":"${auctionHistory.playerId}", "playerName":"${playerName}", "paymentDeadline":"${bid.paymentDeadline}"}`.replace(
            /"/g,
            '\\"'
          ),
        });

        if (destinataryTeamId) {
          await GamelayerService.setMessage({
            destinatary: destinataryTeamId,
            category: 'auction',
            auctionId: auctionHistory.id,
            title: '',
            text: 'auction_seller_sells',
            customImageUrl: '',
            metadata: `{"bidderTeamId":"${bid.teamId}", "bidderTeamName":"${bidderTeamName}", "amount":"${totalAmount}", "playerId":"${auctionHistory.playerId}", "playerName":"${playerName}"}`.replace(
              /"/g,
              '\\"'
            ),
          });
        }
        break;
      case bidStates.accepted:
        if (
          !buyerLosesAuctionTeamsNotified.includes(bid.teamId) &&
          maxBid.teamId != bid.teamId
        ) {
          await GamelayerService.setMessage({
            destinatary: bid.teamId,
            category: 'auction',
            auctionId: auctionHistory.id,
            title: '',
            text: 'auction_buyer_loses_auction',
            customImageUrl: '',
            metadata: `{"amount": "${totalAmount}", "maxBidderTeamId":"${maxBid.teamId}", "maxBidderTeamName":"${maxBidderTeamName}", "playerId": "${auctionHistory.playerId}", "playerName":"${playerName}"}`.replace(
              /"/g,
              '\\"'
            ),
          });
          buyerLosesAuctionTeamsNotified.push(bid.teamId);
        }
        break;
    }
  }
};

module.exports = processPayingAuction;
