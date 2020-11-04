const dayjs = require('dayjs');
const selectLastChecked = require('../../repositories/selectLastChecked');
const updateLastChecked = require('../../repositories/updateLastChecked');
const HorizonService = require('../HorizonService');
const { mailboxCronEntities } = require('../../config');
const logger = require('../../logger');
const processAcceptedBids = require('../bids/processAcceptedBids');

const processBidsHistories = async () => {
  const bidLastChecked = await selectLastChecked({
    entity: mailboxCronEntities.bid,
  });

  const lastBidsHistories = await HorizonService.getLastAcceptedBidsHistories({
    lastChecked: dayjs(bidLastChecked).add(1, 'second').format(),
  });

  const newLastChecked =
    lastBidsHistories[lastBidsHistories.length - 1] &&
    lastBidsHistories[lastBidsHistories.length - 1].insertedAt
      ? lastBidsHistories[lastBidsHistories.length - 1].insertedAt
      : bidLastChecked;

  await updateLastChecked({
    entity: mailboxCronEntities.auction,
    lastChecked: newLastChecked,
  });
  const areNewHistories =
    dayjs(bidLastChecked).format() != dayjs(newLastChecked).format();

  logger.info(
    areNewHistories
      ? `Processing Accepted Bids Histories - LastCheckedAuctionTime: ${bidLastChecked} - NewLastCheckedAuctionTime: ${newLastChecked}`
      : `Processing Accepted Bids Histories - No new accepted bids histories since lastCheckedAuctionTime: ${bidLastChecked}`
  );

  if (areNewHistories) {
    await processAcceptedBids({ bids: lastBidsHistories });
  }
};

module.exports = processBidsHistories;
