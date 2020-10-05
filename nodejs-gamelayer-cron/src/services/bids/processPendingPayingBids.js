const HorizonService = require('../HorizonService');
const dayjs = require('dayjs');
const GamelayerService = require('../Gamelayerservice');
const { timeToNotifyPendinPayingBids } = require('../config');

const processPendingPayingBids = async () => {
  const bids = await HorizonService.getPayingBids();

  for (const bid of bids) {
    const paymentDeadline = dayjs(bid.paymentDeadline);
    const timeRemainingToPay = dayjs().utc().subtract(paymentDeadline).hour();

    if (timeRemainingToPay < timeToNotifyPendinPayingBids) {
      await GamelayerService.setMessage({
        destinatary: bid.teamId,
        category: 'auction',
        auctionId: auctionHistory.id,
        text: 'You should pay the auction, less than 12h remaining',
        customImageUrl: '',
        metadata: '',
      });
    }
  }
};

module.exports = processPendingPayingBids;
