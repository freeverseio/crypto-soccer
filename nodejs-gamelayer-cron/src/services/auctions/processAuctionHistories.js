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
  console.log(
    'processAuctionHistories -> auctionLastChecked',
    auctionLastChecked
  );

  logger.info(
    `Processing Auction Histories - LastCheckedAuctionTime: ${auctionLastChecked}`
  );
  const lastAuctionsHistories = await HorizonService.getLastAuctionsHistories({
    lastChecked: auctionLastChecked,
  });

  const newLastChecked =
    lastAuctionsHistories[lastAuctionsHistories.length - 1].insertedAt;
  await updateLastChecked({
    entity: mailboxTypes.auction,
    lastChecked: newLastChecked,
  });

  logger.info(
    `Processing Auction Histories - LastCheckedAuctionTime: ${auctionLastChecked} - NewLastCheckedAuctionTime: ${newLastChecked}`
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
