const HorizonService = require('../HorizonService');
const utc = require('dayjs/plugin/utc');
const dayjs = require('dayjs');
const GamelayerService = require('../GamelayerService');
const { timeToNotifyPendinPayingBids } = require('../../config');
const logger = require('../../logger');
dayjs.extend(utc);

const processPendingPayingBids = async () => {
  const bids = await HorizonService.getPayingBids();

  logger.info(`Processing Pending Paying Bids`);

  for (const bid of bids) {
    try {
      const paymentDeadline = dayjs.unix(bid.paymentDeadline).utc();
      const timeRemainingToPay = paymentDeadline.diff(dayjs.utc(), 'hour');

      if (timeRemainingToPay < timeToNotifyPendinPayingBids) {
        const messages = await GamelayerService.getMessages({
          teamId: bid.teamId,
          auctionId: bid.auctionId,
        });

        pendingPaymentMessage = messages.find(
          (m) =>
            m.text == 'auction_buyer_pending_payment' &&
            m.auctionId == bid.auctionId
        );

        if (!pendingPaymentMessage) {
          const auction = await HorizonService.getAuction({
            auctionId: bid.auctionId,
          });

          const {
            name: bidderTeamName,
          } = await GamelayerService.getInfoFromTeamId({
            teamId: bid.teamId,
          });
          const totalAmount =
            parseInt(auction.price) + parseInt(bid.extraPrice);

          await GamelayerService.setMessage({
            destinatary: bid.teamId,
            category: 'auction',
            auctionId: bid.auctionId,
            title: '',
            text: 'auction_buyer_pending_payment',
            customImageUrl: '',
            metadata: `{"amount":"${totalAmount}", "playerId":"${auction.playerId}", "playerName":"${auction.playerByPlayerId.name}", "bidderTeamId":"${bid.teamId}", "bidderTeamName":"${bidderTeamName}", "paymentDeadline":"${bid.paymentDeadline}"}`.replace(
              /"/g,
              '\\"'
            ),
          });
        }
      }
    } catch (e) {
      logger.info(
        `Error processing pending paying bid: ${JSON.stringify(bid)}`
      );
      logger.error(e);
    }
  }
};

module.exports = processPendingPayingBids;
