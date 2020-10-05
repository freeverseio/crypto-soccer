const processPayingAuction = require('./processPayingAuction');
const processWithdrawableBySellerAuction = require('./processWithdrawableBySellerAuction');
const selectLastChecked = require('../../repositories/selectLastChecked');
const updateLastChecked = require('../../repositories/updateLastChecked');
const HorizonService = require('../HorizonService');
const { mailboxTypes, auctionStates } = require('../../config');

const processAuctionHistories = async () => {
  const auctionLastChecked = await selectLastChecked({
    entity: mailboxTypes.auction,
  });

  const lastAuctionsHistories = await HorizonService.getLastAuctionsHistories({
    lastChecked: auctionLastChecked,
  });

  const newLastChecked =
    lastAuctionsHistories[lastAuctionsHistories.length - 1].insertedAt;

  await updateLastChecked({
    entity: mailboxTypes.auction,
    lastChecked: newLastChecked,
  });

  for (const auctionHistory of lastAuctionsHistories) {
    switch (auctionHistory.state.toLowerCase()) {
      case auctionStates.withdrawableBySeller:
        await processWithdrawableBySellerAuction({ auctionHistory });
        break;
      case auctionStates.paying:
        await processPayingAuction({ auctionHistory });
        break;
    }
  }
};

module.exports = processAuctionHistories;
