const dayjs = require('dayjs');
const processPayingAuction = require('./processPayingAuction');
const processWithdrawableBySellerAuction = require('./processWithdrawableBySellerAuction');
const selectLastChecked = require('../../repositories/selectLastChecked');
const updateLastChecked = require('../../repositories/updateLastChecked');
const HorizonService = require('../HorizonService');
const { mailboxTypes, auctionStates } = require('../../config');
const logger = require('../../logger');

const processAuctionHistories = async () => {
  const auctionLastChecked = await selectLastChecked({
    entity: mailboxTypes.auction,
  });

  const lastAuctionsHistories = await HorizonService.getLastAuctionsHistories({
    lastChecked: dayjs(auctionLastChecked).format(),
  });

  const newLastChecked =
    lastAuctionsHistories[lastAuctionsHistories.length - 1] &&
    lastAuctionsHistories[lastAuctionsHistories.length - 1].inserted_at
      ? lastAuctionsHistories[lastAuctionsHistories.length - 1].insertedAt
      : auctionLastChecked;

  await updateLastChecked({
    entity: mailboxTypes.auction,
    lastChecked: newLastChecked,
  });

  logger.info(
    auctionLastChecked != newLastChecked
      ? `Processing Auction Histories - LastCheckedAuctionTime: ${auctionLastChecked} - NewLastCheckedAuctionTime: ${newLastChecked}`
      : `Processing Auction Histories - No new auction histories since lastCheckedAuctionTime: ${auctionLastChecked}`
  );

  for (const auctionHistory of lastAuctionsHistories) {
    try {
      switch (auctionHistory.state.toLowerCase()) {
        case auctionStates.withdrawableBySeller:
          await processWithdrawableBySellerAuction({ auctionHistory });
          break;
        case auctionStates.paying:
          await processPayingAuction({ auctionHistory });
          break;
      }
    } catch (e) {
      logger.info(
        `Error processing auction history: ${JSON.stringify(auctionHistory)}`
      );
      logger.error(e);
    }
  }
};

module.exports = processAuctionHistories;
