const dayjs = require('dayjs');
const processPayingAuction = require('./processPayingAuction');
const processWithdrawableBySellerAuction = require('./processWithdrawableBySellerAuction');
const selectLastChecked = require('../../repositories/selectLastChecked');
const updateLastChecked = require('../../repositories/updateLastChecked');
const HorizonService = require('../HorizonService');
const { mailboxTypes, auctionStates } = require('../../config');
const logger = require('../../logger');
const processAcceptedBids = require('../bids/processAcceptedBids');

const processAuctionHistories = async () => {
  const auctionLastChecked = await selectLastChecked({
    entity: mailboxTypes.auction,
  });

  const lastAuctionsHistories = await HorizonService.getLastAuctionsHistories({
    lastChecked: dayjs(auctionLastChecked).add(1, 'second').format(),
  });

  const newLastChecked =
    lastAuctionsHistories[lastAuctionsHistories.length - 1] &&
    lastAuctionsHistories[lastAuctionsHistories.length - 1].insertedAt
      ? lastAuctionsHistories[lastAuctionsHistories.length - 1].insertedAt
      : auctionLastChecked;

  await updateLastChecked({
    entity: mailboxTypes.auction,
    lastChecked: newLastChecked,
  });
  const areNewHistories =
    dayjs(auctionLastChecked).format() != dayjs(newLastChecked).format();

  logger.info(
    areNewHistories
      ? `Processing Auction Histories - LastCheckedAuctionTime: ${auctionLastChecked} - NewLastCheckedAuctionTime: ${newLastChecked}`
      : `Processing Auction Histories - No new auction histories since lastCheckedAuctionTime: ${auctionLastChecked}`
  );

  if (areNewHistories) {
    await processAcceptedBids({
      lastChecked: dayjs(auctionLastChecked).add(1, 'second').format(),
    });

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
  }
};

module.exports = processAuctionHistories;
