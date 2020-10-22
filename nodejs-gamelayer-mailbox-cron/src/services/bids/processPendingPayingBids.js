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
        await GamelayerService.setMessage({
          destinatary: bid.teamId,
          category: 'auction',
          auctionId: bid.auctionId,
          text: 'auction_buyer_pending_payment',
          customImageUrl: '',
          metadata: '{}',
        });
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
