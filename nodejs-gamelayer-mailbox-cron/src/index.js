const cron = require('node-cron');
const { MAILBOX_CRON } = require('./config');
const processOffersHistories = require('./services/offers/processOffersHistories');
const processAuctionHistories = require('./services/auctions/processAuctionHistories');
const processPendingPayingBids = require('./services/bids/processPendingPayingBids');
const processBidsHistories = require('./services/bids/processBidsHistories');
const processUnpayments = require('./services/unpayments/processUnpayments');
const logger = require('./logger');

cron.schedule(MAILBOX_CRON, async () => {
  logger.info('Initiating mailbox events polling...');
  await processOffersHistories();
  await processBidsHistories();
  await processAuctionHistories();
  await processPendingPayingBids();
  await processUnpayments();
});
