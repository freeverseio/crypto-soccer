const cron = require('node-cron');
const { MAILBOX_CRON } = require('./config');
const processOffersHistories = require('./services/offers/processOffersHistories');
const processAuctionHistories = require('./services/auctions/processAuctionHistories');
const processPendingPayingBids = require('./services/bids/processPendingPayingBids');
const logger = require('./logger');

cron.schedule(MAILBOX_CRON, async () => {
  logger.info('Initiating mailbox events polling...');
  await processAuctionHistories();
  await processOffersHistories();
  await processPendingPayingBids();
});
